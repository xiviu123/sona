package sona

import (
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/xiviu123/sona/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server *grpc.Server
}

func NewGRPCServer() *GRPCServer {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			interceptor.RequestIDInterceptor(),
			interceptor.AuthenticationInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	return &GRPCServer{
		server: s,
	}
}

func (s *GRPCServer) Start(addr string) error {
	listenPort, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	reflection.Register(s.server)
	return s.server.Serve(listenPort)
}
