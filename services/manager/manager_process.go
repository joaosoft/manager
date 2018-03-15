package gomanager

// IProcess...
type IProcessManager interface {
	Start() error
	Stop() error
	Started() bool
}
