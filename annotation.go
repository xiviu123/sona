package sona

import (
	"context"
	"net/http"

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

	return metadata.New(map[string]string{
		interceptor.XRequestIDKey: requestID,
	})
}
