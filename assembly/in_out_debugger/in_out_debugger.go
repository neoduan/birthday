package in_out_debugger

import (
	"github.com/neoduan/birthday/assembly"

	"github.com/neoduan/birthday/interceptor/in_out_debugger"
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
		in_out_debugger.UnaryServerInterceptor(),
	)

}

func (this *Assembly) setupSSI() {
	//nothind to do
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
