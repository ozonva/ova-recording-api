package recording

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

const RequestIdKey = "requestId"
const MethodKey = "method"

func RequestIdInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = AddValue(ctx, RequestIdKey, NewRequestId())
	ctx = AddValue(ctx, MethodKey, info.FullMethod)
	h, err := handler(ctx, req)
	return h, err
}

func NewRequestId() uuid.UUID {
	newUuid, err := uuid.NewUUID()
	if err != nil {
		fmt.Printf("cannot get uuid: %s", err)
		return uuid.UUID{}
	}
	return newUuid
}
