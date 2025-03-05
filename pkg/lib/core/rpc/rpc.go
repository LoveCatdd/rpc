package rpc

import (
	"net"

	"github.com/LoveCatdd/util/pkg/lib/core/log"
	"google.golang.org/grpc"
)

func GRPCServer() (*grpc.Server, net.Listener, error) {

	addr := RpcConf.Rpc.Server.Addr
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(arrayInterceptors()), // 注册中间件：日志、auth身份验证、超时
	)

	log.Infof("🚀 gRPC server starting at %s", lis.Addr())
	return grpcServer, lis, nil
}
