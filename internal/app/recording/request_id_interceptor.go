package recording

import (
	"context"
	"google.golang.org/grpc"
)

type RequestID int
const RequestIdKey = "requestId"
const MethodKey = "method"
var currRequestId = 0

func RequestIdInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = AddValue(ctx, RequestIdKey, NewRequestId())
	ctx = AddValue(ctx, MethodKey, info.FullMethod)
	h, err := handler(ctx, req)
	return h, err
}

func NewRequestId() RequestID {
	currRequestId++;
	return RequestID(currRequestId)
}
