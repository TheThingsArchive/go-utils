// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

// Package sendbuffer implements a buffered Client-Streaming RPC that drops the oldest messages on buffer overflow.
package sendbuffer

import (
	"sync/atomic"

	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

var slack = 4

// New returns a new Stream with the given buffer size and setup function.
// If you start calling SendMsg() immediately after this, the stream will start buffering.
// You must call Run() in a separate goroutine to actually start handling the stream.
func New(bufferSize int, setup func() (grpc.ClientStream, error)) *Stream {
	return &Stream{
		setupFunc:  setup,
		sendBuffer: make(chan interface{}, bufferSize+slack),
		log:        ttnlog.Get(),
	}
}

// Stream client->server streaming rpc that buffers (at most) the last {bufferSize} messages
type Stream struct {
	// BEGIN sync/atomic aligned
	sent    uint64
	dropped uint64
	// END sync/atomic aligned

	setupFunc  func() (grpc.ClientStream, error)
	sendBuffer chan interface{}

	log ttnlog.Interface
}

// Stats of the stream
func (s *Stream) Stats() (sent, dropped uint64) {
	return atomic.LoadUint64(&s.sent), atomic.LoadUint64(&s.dropped)
}

// SendMsg sends a message on the stream
func (s *Stream) SendMsg(msg interface{}) {
	if len(s.sendBuffer) > cap(s.sendBuffer)-slack {
		atomic.AddUint64(&s.dropped, 1)
		s.log.Debug("sendbuffer: dropping message")
		<-s.sendBuffer
	}
	s.sendBuffer <- msg
}

// Run the stream
func (s *Stream) Run() error {
	stream, err := s.setupFunc()
	if err != nil {
		return err
	}

	recvErr := make(chan error)
	defer func() {
		go func() { // empty the recvErr channel on return
			<-recvErr
		}()
	}()

	go func() {
		var e empty.Empty
		err := stream.RecvMsg(&e)
		s.log.WithError(err).Debug("sendbuffer: error from stream.RecvMsg")
		recvErr <- err
		close(recvErr)
	}()

	defer stream.CloseSend()

	for {
		select {
		case err := <-recvErr:
			return err
		case <-stream.Context().Done():
			s.log.WithError(stream.Context().Err()).Debug("sendbuffer: context done")
			return stream.Context().Err()
		case msg := <-s.sendBuffer:
			if err = stream.SendMsg(msg); err != nil {
				s.log.WithError(err).Debug("sendbuffer: error from stream.SendMsg")
				return err
			}
			s.log.Debug("sendbuffer: sent message")
			atomic.AddUint64(&s.sent, 1)
		}
	}
}
