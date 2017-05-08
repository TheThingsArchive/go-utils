package rpclog

import (
	"strings"

	"github.com/TheThingsNetwork/go-utils/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type _fields map[string][]string

func (f _fields) add(key string, values ...string) {
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

func (f _fields) LogFields() log.Fields {
	fields := make(log.Fields)
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

func (f _fields) addFromMD(md metadata.MD) {
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

// FieldsFromContext returns peer information and MDLogFields from the given context
func FieldsFromContext(ctx context.Context) log.Fields {
	fields := make(_fields)
	if peer, ok := peer.FromContext(ctx); ok {
		fields.add("caller-ip", peer.Addr.String())
		if peer.AuthInfo != nil {
			fields.add("auth-type", peer.AuthInfo.AuthType())
		}
	}
	if in, ok := metadata.FromIncomingContext(ctx); ok {
		fields.addFromMD(in)
	}
	if out, ok := metadata.FromOutgoingContext(ctx); ok {
		fields.addFromMD(out)
	}
	return fields.LogFields()
}
