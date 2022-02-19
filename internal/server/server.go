package server

import (
	"context"
	"fmt"

	hellov1 "github.com/dtan4/grpc-private-bff-example/api/hello/v1"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
}

var _ hellov1.HelloServiceServer = (*Server)(nil)

func New(logger *zap.Logger) hellov1.HelloServiceServer {
	return &Server{
		logger: logger,
	}
}

func (s *Server) SayHello(ctx context.Context, r *hellov1.SayHelloRequest) (*hellov1.SayHelloResponse, error) {
	s.logger.Info("received SayHello", zap.String("name", r.GetName()))

	return &hellov1.SayHelloResponse{
		Message: fmt.Sprintf("Hello, %s", r.GetName()),
	}, nil
}
