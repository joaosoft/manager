package gomanager

// appConfig ...
type appConfig struct {
	GoManager ManagerConfig `json:"log"`
}

// ManagerConfig ...
type ManagerConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}
