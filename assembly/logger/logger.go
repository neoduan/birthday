package logger

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/neoduan/birthday/assembly"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

var _ assembly.IAssembly = &Assembly{}

type Assembly struct {
	log *zap.Logger
}

func (this *Assembly) Setup() {
	this.setupUSI()
	this.setupSSI()

	this.setupUCI()
	this.setupSCI()
}

func (this *Assembly) Unload() {
	this.log.Sync()
}

func call(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
	return true
}

func (this *Assembly) setupUSI() {
	assembly.USI = append(assembly.USI,
		grpc_zap.UnaryServerInterceptor(this.log),
		//grpc_zap.PayloadUnaryServerInterceptor(this.log, call),
	)

}

func (this *Assembly) setupSSI() {
	assembly.SSI = append(assembly.SSI,
		grpc_zap.StreamServerInterceptor(this.log),
	)
}

func (this *Assembly) setupUCI() {
	assembly.UCI = append(assembly.UCI,
		grpc_zap.UnaryClientInterceptor(this.log),
	)
}

func (this *Assembly) setupSCI() {
	assembly.SCI = append(assembly.SCI,
		grpc_zap.StreamClientInterceptor(this.log),
	)
}

func New(log *zap.Logger) *Assembly {
	return &Assembly{log: log}
}
