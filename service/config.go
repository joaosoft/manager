package gomanager

// appConfig ...
type appConfig struct {
	GoManager goManagerConfig `json:"log"`
}

// goManagerConfig ...
type goManagerConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}
