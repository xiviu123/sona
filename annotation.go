package sona

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/xid"
	"github.com/xiviu123/sona/interceptor"
	"google.golang.org/grpc/metadata"
)

// RequestIDAnnotator takes requestID from http request header and sets it to metadata.
func RequestIDAnnotator(ctx context.Context, req *http.Request) metadata.MD {
	requestID := req.Header.Get(interceptor.XRequestIDKey)
	if requestID == "" {
		requestID = xid.New().String()
	}

	md := make(map[string]string)
	md[interceptor.XRequestIDKey] = requestID

	if method, ok := runtime.RPCMethod(ctx); ok {
		md["method"] = method // /grpc.gateway.examples.internal.proto.examplepb.LoginService/Login
	}

	if pattern, ok := runtime.HTTPPathPattern(ctx); ok {
		md["pattern"] = pattern // /v1/example/login
	}

	return metadata.New(md)
}
