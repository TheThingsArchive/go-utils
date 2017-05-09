// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package auth

import (
	"github.com/TheThingsNetwork/go-utils/grpc/ttnctx"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const tokenKey = "token"

// TokenCredentials RPC Credentials
type TokenCredentials struct {
	allowInsecure bool
	token         string
	tokenFunc     func(id string) string
	tokenFuncKey  string
}

// WithInsecure returns a copy of the TokenCredentials, allowing insecure transport
func (c *TokenCredentials) WithInsecure() *TokenCredentials {
	return &TokenCredentials{token: c.token, tokenFunc: c.tokenFunc, allowInsecure: true}
}

// WithStaticToken injects a static token on each request
func WithStaticToken(token string) *TokenCredentials {
	return &TokenCredentials{
		token: token,
	}
}

// WithTokenFunc returns TokenCredentials that execute the tokenFunc on each request
// The value of v sent to the tokenFunk is the MD value of the supplied k
func WithTokenFunc(k string, tokenFunc func(v string) string) *TokenCredentials {
	return &TokenCredentials{
		tokenFunc:    tokenFunc,
		tokenFuncKey: k,
	}
}

// RequireTransportSecurity implements credentials.PerRPCCredentials
func (c *TokenCredentials) RequireTransportSecurity() bool { return !c.allowInsecure }

// GetRequestMetadata implements credentials.PerRPCCredentials
func (c *TokenCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	md := ttnctx.MetadataFromOutgoingContext(ctx)
	token, _ := ttnctx.TokenFromMetadata(md)
	if token != "" {
		return map[string]string{tokenKey: token}, nil
	}
	if c.tokenFunc != nil {
		if v, ok := md[c.tokenFuncKey]; ok && len(v) > 0 {
			return map[string]string{tokenKey: c.tokenFunc(v[0])}, nil
		}
	}
	if c.token != "" {
		return map[string]string{tokenKey: c.token}, nil
	}
	return map[string]string{tokenKey: ""}, nil
}

// DialOption returns a DialOption for the TokenCredentials
func (c *TokenCredentials) DialOption() grpc.DialOption {
	return grpc.WithPerRPCCredentials(c)
}
