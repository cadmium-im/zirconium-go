package core

const (
	ModuleInterfaceName = "Module"
)

type Module interface {
	Initialize(moduleAPI ModuleAPI)
	Name() string
	Version() string
}

type ModuleRef struct {
	F func() Module
}

type ModuleFunc func() Module
