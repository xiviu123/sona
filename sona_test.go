package sona

import (
	"context"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xiviu123/sona/test"
)

func TestGateway(t *testing.T) {
	server := NewGateway()

	grpcServerEndpoint := flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
	grpcServerEndpoint1 := flag.String("grpc-server-endpoint1", "localhost:9091", "gRPC server endpoint")
	server.AddServiceHandle(grpcServerEndpoint, test.RegisterExampleServiceHandlerFromEndpoint)
	server.AddServiceHandle(grpcServerEndpoint1, test.RegisterExampleService2HandlerFromEndpoint)
	err := server.Start(":8080")

	assert.Equal(t, nil, err)
}

type exampleServiceServer struct {
}

func (s *exampleServiceServer) Ping(ctx context.Context, in *test.PingRequest) (*test.PingResponse, error) {
	return &test.PingResponse{Msg: "hehe"}, nil
}

func NewExampleServiceServer() test.ExampleServiceServer {
	return &exampleServiceServer{}
}

func TestService(t *testing.T) {
	s := NewGRPCServer()

	test.RegisterExampleServiceServer(s.server, NewExampleServiceServer())
	s.Start(":9090")
}
