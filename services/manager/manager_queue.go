package gomanager

// IQueue ...
type IQueue interface {
	Start() error
	Stop() error
}
