package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	hellov1 "github.com/dtan4/grpc-private-bff-example/api/hello/v1"
	"github.com/dtan4/grpc-private-bff-example/internal/log"
	"github.com/dtan4/grpc-private-bff-example/internal/server"
)

const (
	defaultPort = 8080
)

func realMain() error {
	port := defaultPort
	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("parse port number: %w", err)
		}
		port = p
	}

	logger, err := log.NewLogger()
	if err != nil {
		return fmt.Errorf("create new logger: %w", err)
	}
	defer logger.Sync()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("listen port %d: %w", port, err)
	}
	defer l.Close()

	s := grpc.NewServer()
	hellov1.RegisterHelloServiceServer(s, server.New(logger))

	logger.Info("starting server", zap.Int("port", port))

	if err := s.Serve(l); err != nil {
		return fmt.Errorf("run gRPC server: %w", err)
	}

	return nil
}

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
