// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package restartstream

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/ttn/utils/backoff"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
)

// Settings for Interceptor
type Settings struct {
	RetryableCodes []codes.Code
	Backoff        backoff.Config
}

// DefaultSettings for Interceptor
var DefaultSettings = Settings{
	RetryableCodes: []codes.Code{
		codes.Canceled, // context.WithCancel
		codes.Unknown,
		codes.DeadlineExceeded, // context.WithDeadline
		codes.Aborted,
		codes.Unavailable,
		codes.Internal,
	},
	Backoff: backoff.DefaultConfig,
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
	backoff        backoff.Config
	retries        int

	sync.RWMutex
	grpc.ClientStream
}

func (s *restartingStream) start() (err error) {
	s.Lock()
	defer s.Unlock()

stream:
	for {
		s.log.Debug("Stream (re)starting")
		s.ClientStream, err = s.streamer(s.ctx, s.desc, s.cc, s.method, s.opts...)
		if err == nil {
			s.retries = 0
			break
		}
		s.log.WithField("error", grpc.ErrorDesc(err)).Debug("Stream setup unsuccessful")
		for _, retryable := range s.retryableCodes {
			if grpc.Code(err) == retryable {
				backoff := s.backoff.Backoff(s.retries)
				s.log.WithField("Duration", backoff).Debug("Stream backing off")
				time.Sleep(backoff)
				s.retries++
				continue stream
			}
		}
		return err
	}
	log := s.log
	if peer, ok := peer.FromContext(s.ClientStream.Context()); ok {
		log = log.WithField("Peer", peer.Addr)
	}
	log.Debug("Stream started")

	return
}

// ErrStreamClosed is returned when trying to call SendMsg or RecvMsg on a closed stream
var ErrStreamClosed = errors.New("grpc: stream closed")

func (s *restartingStream) SendMsg(m interface{}) error {
	s.RLock()
	stream := s.ClientStream
	s.RUnlock()
	if stream == nil {
		return ErrStreamClosed
	}

	return stream.SendMsg(m) // blocking
}

func (s *restartingStream) RecvMsg(m interface{}) error {
	s.RLock()
	stream := s.ClientStream
	s.RUnlock()
	if stream == nil {
		return ErrStreamClosed
	}

	err := stream.RecvMsg(m) // blocking
	if err == nil {
		return nil
	}

	stream.CloseSend()

	if err := s.ctx.Err(); err != nil {
		return err // context canceled
	}

	if err == io.EOF {
		return err // eof
	}

	for _, retryable := range s.retryableCodes {
		if grpc.Code(err) == retryable {
			if err := s.start(); err != nil {
				return err
			}
			return s.RecvMsg(m)
		}
	}

	return err
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
		s := &restartingStream{
			log: log.Get().WithField("Method", method),

			ctx:      ctx,
			desc:     desc,
			cc:       cc,
			method:   method,
			streamer: streamer,
			opts:     opts,

			retryableCodes: settings.RetryableCodes,
			backoff:        settings.Backoff,
		}

		err = s.start()

		stream = s
		return
	}
}
