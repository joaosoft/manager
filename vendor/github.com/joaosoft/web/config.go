package web

type AppServerConfig struct {
	Server ServerConfig `json:"server"`
}

type AppClientConfig struct {
	Client ClientConfig `json:"client"`
}

type Log struct {
	Level string `json:"level"`
}
