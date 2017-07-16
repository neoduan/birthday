package assembly

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var (
	USI []grpc.UnaryServerInterceptor  = make([]grpc.UnaryServerInterceptor, 0)
	SSI []grpc.StreamServerInterceptor = make([]grpc.StreamServerInterceptor, 0)
	UCI []grpc.UnaryClientInterceptor  = make([]grpc.UnaryClientInterceptor, 0)
	SCI []grpc.StreamClientInterceptor = make([]grpc.StreamClientInterceptor, 0)
)

func WithUnartServer() grpc.ServerOption {
	return grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(USI...))
}

func WithStreamServer() grpc.ServerOption {
	return grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(SSI...))
}

func WithStreamClient() grpc.DialOption {
	return grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(UCI...))
}

func WithUnartClient() grpc.DialOption {
	return grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(SCI...))
}

type IAssembly interface {
	Setup()
	Unload()
}

func Setup(assemblys ...IAssembly) {
	for _, assembly := range assemblys {
		assembly.Setup()
	}
}

func Unload(assemblys ...IAssembly) {
	for _, assembly := range assemblys {
		assembly.Unload()
	}
}
