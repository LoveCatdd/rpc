package rpc

import (
	"context"
	"errors"

	gcontext "github.com/LoveCatdd/webctx/pkg/lib/core/context"

	"github.com/LoveCatdd/util/pkg/lib/core/log"
	"github.com/LoveCatdd/webctx/pkg/lib/core/goroutine"
	"github.com/LoveCatdd/webctx/pkg/lib/core/web/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor 认证拦截器
func (InterceptorImpl) AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Info("Missing metadata")
			return nil, errors.New("missing metadata")
		}

		token := md.Get(auth.JWT_AUTHORIZATION_KEY)
		if len(token) == 0 {
			log.Info("Missing token")
			return nil, errors.New("missing token")
		}

		mapClaims, err := auth.ExtractMapClaims(token[0])
		if err != nil {
			log.Info(" Unauthorized request")
			return nil, errors.New("unauthorized")
		}

		contextHolder := new(goroutine.GoroutineContextHolder)
		contextHolder.Initialization()

		// mapClaims 维护到 contextHolder
		contextHolder.Change(goroutine.JWT_MAP_CLAIM, mapClaims)

		// 持久化到 request_context 中
		customContext := gcontext.NewCustomContext(contextHolder)

		//维护到ctx中
		ctx = metadata.NewIncomingContext(ctx, md)
		ctx = context.WithValue(ctx, gcontext.CustonContextKey, customContext)

		return handler(ctx, req)
	}
}
