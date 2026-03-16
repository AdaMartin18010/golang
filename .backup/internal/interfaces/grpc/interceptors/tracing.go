package interceptors

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TracingUnaryInterceptor 追踪拦截器
func TracingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	tracer := otel.Tracer("grpc")
	ctx, span := tracer.Start(ctx, info.FullMethod,
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	// 从 metadata 中提取追踪信息
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// 可以提取 TraceContext 等
		_ = md
	}

	resp, err := handler(ctx, req)

	if err != nil {
		s, _ := status.FromError(err)
		span.SetStatus(codes.Error, s.Message())
		span.RecordError(err)
	} else {
		span.SetStatus(codes.Ok, "success")
	}

	span.SetAttributes(
		attribute.String("grpc.method", info.FullMethod),
		attribute.String("grpc.status", status.Code(err).String()),
	)

	return resp, err
}
