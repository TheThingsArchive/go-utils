// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package rpclog

import (
	"strings"

	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type fieldMap map[string][]string

func (f fieldMap) add(key string, values ...string) {
	if _, ok := f[key]; !ok {
		f[key] = make([]string, 0, len(values))
	}
	for _, value := range values {
		if value == "" {
			continue
		}
		f[key] = append(f[key], value)
	}
}

func (f fieldMap) LogFields() ttnlog.Fields {
	fields := make(ttnlog.Fields)
	for k, v := range f {
		switch len(v) {
		case 0:
		case 1:
			fields[k] = v[0]
		default:
			fields[k] = strings.Join(v, ",")
		}
	}
	return fields
}

// MDLogFields are logged from the context
var MDLogFields = []string{"id", "service-name", "service-version", "limit", "offset"}

func (f fieldMap) addFromMD(md metadata.MD) {
	for _, key := range MDLogFields {
		if v, ok := md[key]; ok {
			f.add(key, v...)
		}
	}
	if v, ok := md["key"]; ok && len(v) > 0 && v[0] != "" {
		f.add("auth-type", "key")
	}
	if v, ok := md["token"]; ok && len(v) > 0 && v[0] != "" {
		f.add("auth-type", "token")
	}
}

func (f fieldMap) addFromPeer(peer *peer.Peer) {
	f.add("caller-ip", peer.Addr.String())
	if peer.AuthInfo != nil {
		f.add("auth-type", peer.AuthInfo.AuthType())
	}
}

// FieldsFromIncomingContext returns peer information and MDLogFields from the given context
func FieldsFromIncomingContext(ctx context.Context) ttnlog.Fields {
	fields := make(fieldMap)
	if peer, ok := peer.FromContext(ctx); ok {
		fields.addFromPeer(peer)
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fields.addFromMD(md)
	}
	return fields.LogFields()
}

// FieldsFromOutgoingContext returns peer information and MDLogFields from the given context
func FieldsFromOutgoingContext(ctx context.Context) ttnlog.Fields {
	fields := make(fieldMap)
	if peer, ok := peer.FromContext(ctx); ok {
		fields.addFromPeer(peer)
	}
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		fields.addFromMD(md)
	}
	return fields.LogFields()
}
