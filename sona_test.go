package sona

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xiviu123/sona/lock"
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

	test.RegisterExampleServiceServer(s.Server(), NewExampleServiceServer())
	s.Start(":9090")
}

type RedisCfg struct {
	host string
	pwd  string
}

func (r *RedisCfg) Host() string {
	return r.host
}

func (r *RedisCfg) Pwd() string {
	return r.pwd
}
func TestLock(t *testing.T) {

	cfg := &RedisCfg{
		host: "localhost:6379",
		pwd:  "",
	}

	remote := lock.NewRedisStorage(cfg)

	lock, _ := remote.ObtainLock("lockkey", 10*time.Minute)

	_, err := remote.ObtainLock("lockkey", 10*time.Minute)

	assert.EqualError(t, err, err.Error())

	defer func() {
		if lock != nil {
			_ = lock.Unlock()
		}
	}()

}
