package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	hellov1 "github.com/dtan4/grpc-private-bff-example/api/hello/v1"
	"github.com/dtan4/grpc-private-bff-example/internal/log"
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

	crEndpoint := os.Getenv("CLOUD_RUN_ENDPOINT")
	impersonateSAEmail := os.Getenv("IMPERSONATE_SA_EMAIL")

	logger, err := log.NewLogger()
	if err != nil {
		return fmt.Errorf("create new logger: %w", err)
	}
	defer logger.Sync()

	ctx := context.Background()

	mux := runtime.NewServeMux()

	creds, err := makeCreds(ctx, fmt.Sprintf("https://%s", crEndpoint), impersonateSAEmail)
	if err != nil {
		return fmt.Errorf("make token: %w", err)
	}

	// https://cloud.google.com/run/docs/triggering/grpc#connect
	srs, err := x509.SystemCertPool()
	if err != nil {
		return fmt.Errorf("get system cert pool: %w", err)
	}
	tcreds := credentials.NewTLS(&tls.Config{
		RootCAs: srs,
	})

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tcreds),
		grpc.WithPerRPCCredentials(creds),
	}

	if err := hellov1.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:443", crEndpoint), opts); err != nil {
		return fmt.Errorf("register gRPC gateway endpoint: %w", err)
	}

	logger.Info("starting server", zap.Int("port", port))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
