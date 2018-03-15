package gomanager

type IWorkManager interface {
	Start() error
	Stop() error
}
