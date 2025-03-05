package rpc

import (
	"context"
	"time"

	"github.com/LoveCatdd/util/pkg/lib/core/ids"
	"github.com/LoveCatdd/util/pkg/lib/core/log"
	"google.golang.org/grpc"
)

// LoggingInterceptor 记录 gRPC 请求和响应日志
func (InterceptorImpl) LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		log.SetTraceId(ids.UUIDV1())
		log.
			Infof("GRPC request start: , method: %v, request_body: %v",
				info.FullMethod, req)

		resp, err = handler(ctx, req)
		cost := time.Since(start)

		log.
			Infof("GRPC request end:  method: %v, cost: %v, response_body: %v, error: %v",
				info.FullMethod, cost, resp, err)

		return
	}

}
