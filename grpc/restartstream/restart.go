// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package restartstream

import (
	"io"
	"sync"

	"github.com/TheThingsNetwork/go-utils/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
)

// Settings for Interceptor
type Settings struct {
	RetryableCodes []codes.Code
}

// DefaultSettings for Interceptor
var DefaultSettings = Settings{
	RetryableCodes: []codes.Code{
		codes.Canceled,
		codes.Unknown,
		codes.Aborted,
		codes.Unavailable,
	},
}

type restartingStream struct {
	log log.Interface

	ctx      context.Context
	desc     *grpc.StreamDesc
	cc       *grpc.ClientConn
	method   string
	streamer grpc.Streamer
	opts     []grpc.CallOption

	retryableCodes []codes.Code

	done chan struct{}

	sync.RWMutex
	grpc.ClientStream
	closing bool
}

func (s *restartingStream) start() (err error) {
	s.Lock()
	defer s.Unlock()
	s.ClientStream, err = s.streamer(s.ctx, s.desc, s.cc, s.method, s.opts...)
	if err != nil {
		return
	}
	log := s.log
	if peer, ok := peer.FromContext(s.ClientStream.Context()); ok {
		log = log.WithField("Peer", peer.Addr)
	}
	log.Debug("Stream (re)starting")
	return
}

func (s *restartingStream) SendMsg(m interface{}) error {
	s.RLock()
	stream := s.ClientStream
	s.RUnlock()
	return stream.SendMsg(m) // blocking
}

func (s *restartingStream) RecvMsg(m interface{}) error {
	s.RLock()
	stream := s.ClientStream
	s.RUnlock()

	err := stream.RecvMsg(m) // blocking
	if err == nil {
		return nil
	}

	stream.CloseSend()

	s.RLock()
	closing := s.closing
	s.RUnlock()

	if closing {
		return err
	}

	if s.desc.ServerStreams && err == io.EOF {
		close(s.done)
		return err
	}

	for _, retryable := range s.retryableCodes {
		if grpc.Code(err) == retryable {
			s.start()
			return err
		}
	}

	return err
}

func (s *restartingStream) CloseSend() error {
	if s.desc.ClientStreams {
		close(s.done)
	}
	return s.ClientStream.CloseSend()
}

// Interceptor automatically restarts streams on non-expected errors
// To do so, the application should create a for-loop around RecvMsg, which
// returns the same errors that are received from the server.
//
// An io.EOF indicates the end of the stream
//
// To stop the reconnect behaviour, you have to cancel the context
func Interceptor(settings Settings) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (stream grpc.ClientStream, err error) {
		s := new(restartingStream)
		s.retryableCodes = settings.RetryableCodes
		s.ctx = ctx
		s.desc = desc
		s.cc = cc
		s.method = method
		s.streamer = streamer
		s.opts = opts
		s.done = make(chan struct{})

		s.log = log.Get().WithField("Method", method)

		go func() {
			select {
			case <-ctx.Done(): // canceled
			case <-s.done: // eof
			}
			s.Lock()
			defer s.Unlock()
			s.closing = true
		}()

		err = s.start()

		stream = s
		return
	}
}
