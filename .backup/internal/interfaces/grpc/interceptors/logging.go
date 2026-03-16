package interceptors

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

// LoggingUnaryInterceptor 日志拦截器
func LoggingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	slog.Info("gRPC request",
		"method", info.FullMethod,
		"request", req,
	)

	resp, err := handler(ctx, req)

	if err != nil {
		slog.Error("gRPC error",
			"method", info.FullMethod,
			"error", err,
		)
	} else {
		slog.Info("gRPC response",
			"method", info.FullMethod,
		)
	}

	return resp, err
}
