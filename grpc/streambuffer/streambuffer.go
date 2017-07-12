// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

// Package streambuffer implements a buffered streaming RPC that drops the oldest messages on buffer overflow.
package streambuffer

import (
	"context"
	"sync"
	"sync/atomic"

	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// New returns a new Stream with the given buffer size and setup function.
func New(bufferSize int, setup func() (grpc.ClientStream, error)) *Stream {
	return &Stream{
		setupFunc:  setup,
		sendBuffer: make(chan interface{}, bufferSize),
		log:        ttnlog.Get(),
	}
}

// Stream implements a buffered overlay on top of a streaming RPC.
//
// If the buffer is full, the oldest items in the buffer will be dropped.
//
// - Create a new Stream with the New() func.
// - You must call Run() in a separate goroutine to actually start handling the stream. Run calls the setup func you provided in New().
// - The goroutine that calls Run() is responsible for handling backoff.
// - You can start calling SendMsg() immediately after New(), the stream will start buffering until the stream is started by Run().
// - If you want to receive on the stream, Recv() must be called after New(), but before Run().
type Stream struct {
	// BEGIN sync/atomic aligned
	sent     uint64
	received uint64
	dropped  uint64
	// END sync/atomic aligned

	mu sync.RWMutex // Lock while stream is running

	setupFunc  func() (grpc.ClientStream, error)
	recvFunc   func() interface{}
	sendBuffer chan interface{}
	recvBuffer chan interface{}

	log ttnlog.Interface
}

// SetLogger sets the logger for this streambuffer
func (s *Stream) SetLogger(log ttnlog.Interface) {
	s.mu.Lock()
	s.log = log
	s.mu.Unlock()
}

// Recv returns a buffered channel (of the size given to New) that receives messages from the stream.
// The given recv func should return a new proto of the type that you want to receive
// If you want to receive, Recv() must be called BEFORE Run()
func (s *Stream) Recv(recv func() interface{}) <-chan interface{} {
	s.mu.Lock()
	s.recvFunc = recv
	buf := make(chan interface{}, cap(s.sendBuffer))
	s.recvBuffer = buf
	s.mu.Unlock()
	return buf
}

// CloseRecv closes the receive channel
func (s *Stream) CloseRecv() {
	s.mu.Lock()
	if s.recvBuffer != nil {
		close(s.recvBuffer)
		s.recvBuffer = nil
	}
	s.mu.Unlock()
}

// Stats of the stream
func (s *Stream) Stats() (sent, dropped uint64) {
	return atomic.LoadUint64(&s.sent), atomic.LoadUint64(&s.dropped)
}

// SendMsg sends a message (possibly dropping a message on full buffers)
func (s *Stream) SendMsg(msg interface{}) {
	select {
	case s.sendBuffer <- msg: // normal flow if the channel is not blocked
	default:
		s.log.Debug("streambuffer: dropping message before send")
		atomic.AddUint64(&s.dropped, 1)
		<-s.sendBuffer // drop oldest and try again (if conn temporarily unavailable)
		select {
		case s.sendBuffer <- msg:
		default: // drop newest (too many cuncurrent SendMsg)
			s.log.Debug("streambuffer: dropping message before send")
			atomic.AddUint64(&s.dropped, 1)
		}
	}
}

// recvMsg receives a message (possibly dropping a message on full buffers)
func (s *Stream) recvMsg(msg interface{}) {
	s.mu.RLock()
	if s.recvBuffer == nil {
		s.mu.RUnlock()
		return
	}
	defer s.mu.RUnlock()
	select {
	case s.recvBuffer <- msg: // normal flow if the channel is not blocked
	default:
		s.log.Debug("streambuffer: dropping received message")
		atomic.AddUint64(&s.dropped, 1)
		<-s.recvBuffer // drop oldest and try again (if application temporarily unavailable)
		select {
		case s.recvBuffer <- msg:
		default: // drop newest (too many cuncurrent recvMsg)
			atomic.AddUint64(&s.dropped, 1)
			s.log.Debug("streambuffer: dropping received message")
		}
	}
}

// Run the stream.
//
// This calls the underlying grpc.ClientStreams methods to send and receive messages over the stream.
// Run returns the error returned by any of those functions, or context.Canceled if the context is canceled.
func (s *Stream) Run() (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	defer func() {
		if err != nil {
			if grpc.Code(err) == codes.Canceled {
				s.log.Debug("streambuffer: context canceled")
				err = context.Canceled
				return
			}
			if grpc.Code(err) == codes.DeadlineExceeded {
				s.log.Debug("streambuffer: context deadline exceeded")
				err = context.DeadlineExceeded
				return
			}
		}
	}()

	stream, err := s.setupFunc()
	if err != nil {
		s.log.WithError(err).Debug("streambuffer: setup returned error")
		return err
	}

	recvErr := make(chan error)
	defer func() {
		go func() { // empty the recvErr channel on return
			<-recvErr
		}()
	}()

	go func() {
		for {
			var r interface{}
			if s.recvFunc != nil {
				r = s.recvFunc()
			} else {
				r = new(empty.Empty) // Generic proto message if not interested in received values
			}
			err := stream.RecvMsg(r)
			if err != nil {
				s.log.WithError(err).Debug("streambuffer: error from stream.RecvMsg")
				recvErr <- err
				close(recvErr)
				return
			}
			if s.recvFunc != nil {
				s.recvMsg(r)
			}
		}
	}()

	defer stream.CloseSend()

	for {
		select {
		case err := <-recvErr:
			return err
		case <-stream.Context().Done():
			s.log.WithError(stream.Context().Err()).Debug("streambuffer: context done")
			return stream.Context().Err()
		case msg := <-s.sendBuffer:
			if err = stream.SendMsg(msg); err != nil {
				s.log.WithError(err).Debug("streambuffer: error from stream.SendMsg")
				return err
			}
			atomic.AddUint64(&s.sent, 1)
		}
	}
}
