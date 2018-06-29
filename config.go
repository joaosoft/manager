package manager

// AppConfig ...
type AppConfig struct {
	manager ManagerConfig `json:"manager"`
}

// ManagerConfig ...
type ManagerConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"logger"`
}
