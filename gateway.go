package sona

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"github.com/xiviu123/sona/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

type RegisterGatewayEndpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

type Gateway struct {
	handlerMap map[*string]RegisterGatewayEndpoint
}

func NewGateway() *Gateway {
	return &Gateway{
		handlerMap: make(map[*string]RegisterGatewayEndpoint),
	}
}

func grpcDialOptions() []grpc.DialOption {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	decider := func(ctx context.Context, fullMethodName string) bool {
		return true
	}
	startTimeFunc := func() time.Time {
		return time.Now()
	}
	durationFunc := func(startTime time.Time) time.Duration {
		return time.Now().Sub(startTime)
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		// Output request/response payload logs
		grpc.WithUnaryInterceptor(interceptor.PayloadUnaryClientInterceptor(logrus.NewEntry(l), decider, startTimeFunc, durationFunc)),
	}

	return opts
}

func (s *Gateway) AddServiceHandle(endpoint *string, serviceHandle RegisterGatewayEndpoint) {
	s.handlerMap[endpoint] = serviceHandle
}

func (s *Gateway) registerServiceHandlers(ctx context.Context, mux *runtime.ServeMux) error {
	opts := grpcDialOptions()
	for endpoint, handle := range s.handlerMap {
		if err := handle(ctx, mux, *endpoint, opts); err != nil {
			return err
		}
	}

	return nil
}

func (s *Gateway) Start(addr string) error {
	if err := s.run(addr); err != nil {
		return err
	}
	return nil
}

func (s *Gateway) run(addr string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		// Set request_id to grpc metadata
		runtime.WithMetadata(RequestIDAnnotator),
		// runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: false, EmitDefaults: true}),
	)

	if err := s.registerServiceHandlers(ctx, mux); err != nil {
		return err
	}

	handler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Accept-Encoding", "Accept"}),
	)(mux)

	fmt.Printf("http server started on %s\n", addr)
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(addr, handler)
}
