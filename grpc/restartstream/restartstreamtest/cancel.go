// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package restartstreamtest

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Cancel struct {
	cancel chan struct{}
}

func NewCancel() *Cancel {
	return &Cancel{
		cancel: make(chan struct{}),
	}
}

func (c *Cancel) Cancel() {
	close(c.cancel)
	c.cancel = make(chan struct{})
}

func (c *Cancel) Interceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (stream grpc.ClientStream, err error) {
	ctx, cancel := context.WithCancel(ctx)
	stream, err = streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	go func() {
		select {
		case <-stream.Context().Done():
			return
		case <-c.cancel:
			cancel()
		}
	}()
	return
}
