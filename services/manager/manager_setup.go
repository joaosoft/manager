package gomanager

// ISetupManager ...
type ISetupManager interface {
	Get(key string) interface{}
	Set(key string, value string) interface{}
	Reload() error
}
