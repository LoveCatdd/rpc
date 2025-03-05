package rpc

import (
	"context"
	"reflect"

	"github.com/LoveCatdd/util/pkg/lib/core/log"
	"github.com/LoveCatdd/util/pkg/lib/core/viper"
	"github.com/LoveCatdd/webctx/pkg/lib/core/web/auth"
	"google.golang.org/grpc"
)

type Interceptor interface {
	LoggingInterceptor() grpc.UnaryServerInterceptor
	AuthInterceptor() grpc.UnaryServerInterceptor
	TimeoutInterceptor() grpc.UnaryServerInterceptor
}

type InterceptorImpl struct{}

var (
	interceptors map[string]grpc.UnaryServerInterceptor
	ns           []string
)

func init() {
	viper.Yaml(auth.JwtConfig)

	viper.Yaml(RpcConf)

	viper.Yaml(log.Config)
	if log.Config.Zap.Enable { // 开启

		log.InitZap()

		defer log.Sync()
	}

	ns = make([]string, 0)
	interceptors = make(map[string]grpc.UnaryServerInterceptor)

	// 解析yaml
	v := reflect.ValueOf(InterceptorImpl{})

	for _, handler := range RpcConf.Rpc.Handler {
		ns = append(ns, handler.Type)
		method := v.MethodByName(handler.Method)

		if method.IsValid() {

			// 断言成 grpc.UnaryServerInterceptor
			interceptorFunc, ok := method.Interface().(func() grpc.UnaryServerInterceptor)
			if !ok {
				log.Fatalf("interceptorFunc register failed at: %v\n", handler.Type)
				return
			}
			interceptors[handler.Type] = interceptorFunc()
		}
	}

	log.Infof("interceptorFunc register successful!, interceptors:%v", interceptors)
}

// ChainUnaryInterceptors 组合多个拦截器
func chainUnaryInterceptors(interceptors map[string]grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 递归执行拦截器
		chain := handler
		for i := len(ns) - 1; i > -1; i-- {
			ik := ns[i]
			if !checkList(info.FullMethod, ik) {
				chain = wrapUnaryInterceptor(interceptors[ik], info, chain)
			}
		}

		return chain(ctx, req)
	}

}

func wrapUnaryInterceptor(interceptor grpc.UnaryServerInterceptor, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {

		return interceptor(ctx, req, info, next)
	}
}

func arrayInterceptors() grpc.UnaryServerInterceptor {
	return chainUnaryInterceptors(interceptors)
}

func checkList(path, ik string) bool {
	// 匹配白名单
	paths := RpcConf.Rpc.Unauth.Path
	for _, str := range paths {
		if path == str {
			return true && ik == "auth"
		}
	}
	return false
}
