package sona

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/xiviu123/sona/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

type RegisterGatewayEndpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

type Gateway struct {
	handlerMap        map[*string]RegisterGatewayEndpoint
	optionsMap        []runtime.ServeMuxOption
	domainOrigins     []string
	gatewayMiddleware func(h http.Handler) http.Handler
}

func NewGateway() *Gateway {
	return &Gateway{
		handlerMap: make(map[*string]RegisterGatewayEndpoint),
		optionsMap: []runtime.ServeMuxOption{
			runtime.WithMetadata(RequestIDAnnotator),
		},
		domainOrigins:     []string{"*"},
		gatewayMiddleware: nil,
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

func (s *Gateway) AddServiceHandle(endpoint *string, serviceHandle RegisterGatewayEndpoint) *Gateway {
	s.handlerMap[endpoint] = serviceHandle

	return s
}

func (s *Gateway) WithGatewayMidleware(f func(h http.Handler) http.Handler) *Gateway {
	s.gatewayMiddleware = f
	return s
}

func (s *Gateway) WithServeMuxOption(opt runtime.ServeMuxOption) *Gateway {
	s.optionsMap = append(s.optionsMap, opt)
	return s
}

func (s *Gateway) WithCORSOrigins(domains []string) *Gateway {
	s.domainOrigins = domains
	return s
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

	mux := runtime.NewServeMux(s.optionsMap...)

	if err := s.registerServiceHandlers(ctx, mux); err != nil {
		return err
	}

	handler := handlers.CORS(
		handlers.AllowedOrigins(s.domainOrigins),
		handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Accept-Encoding", "Accept"}),
	)(mux)

	hdl := handler
	if s.gatewayMiddleware != nil {
		hdl = s.gatewayMiddleware(handler)
	}

	fmt.Printf("http server started on %s\n", addr)
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(addr, hdl)
}
