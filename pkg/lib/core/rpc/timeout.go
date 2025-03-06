package rpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func (InterceptorImpl) TimeoutInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {

		ctx, cancel := context.WithTimeout(ctx, time.Duration(RpcConf.Rpc.Timeout)*time.Millisecond)
		defer cancel()

		select {
		case <-ctx.Done():
			return nil, status.Error(status.Code(ctx.Err()), "request timeout")

		default:
			return handler(ctx, req)
		}
	}
}
