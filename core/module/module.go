package module

type ModuleInterface interface {
	Name() string
	Init() bool
	Update() bool
	End() bool
}

type ModuleEntity struct {
	signal   chan bool
	interval int64
	lastTs   int64
	module   ModuleInterface
}
