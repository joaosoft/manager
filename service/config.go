package gomanager

// appConfig ...
type appConfig struct {
	GoManager ManagerConfig `json:"gomanager"`
}

// ManagerConfig ...
type ManagerConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}
