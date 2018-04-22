package gomanager

// appConfig ...
type appConfig struct {
	GoManager GoManagerConfig `json:"log"`
}

// GoManagerConfig ...
type GoManagerConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}
