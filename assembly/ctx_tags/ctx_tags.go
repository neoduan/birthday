package ctx_tags

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/neoduan/birthday/assembly"
)

var _ assembly.IAssembly = &Assembly{}

type Assembly struct{}

func (this *Assembly) Setup() {
	this.setupUSI()
	this.setupSSI()

	this.setupUCI()
	this.setupSCI()
}

func (this *Assembly) Unload() {
	//nothing to do
}

func (this *Assembly) setupUSI() {
	assembly.USI = append(assembly.USI,
		grpc_ctxtags.UnaryServerInterceptor(),
	)

}

func (this *Assembly) setupSSI() {
	assembly.SSI = append(assembly.SSI,
		grpc_ctxtags.StreamServerInterceptor(),
	)
}

func (this *Assembly) setupUCI() {
	//nothing to do
}

func (this *Assembly) setupSCI() {
	//nothing to do
}

func New() *Assembly {
	return &Assembly{}
}
