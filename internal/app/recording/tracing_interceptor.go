package recording

import (
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)
import "context"

type SpanKeyType string

const (
	SpanKey SpanKeyType = "opentracing.Span"
)

func TracingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
	h, err := handler(ctx, req)
	span.Finish()
	return h, err
}
