package runner

import (
	"context"

	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

// WebConfig...
type WebConfig struct {
	Address string `json:"address"`
}

// NewWebConfig...
func NewWebConfig(address string) *WebConfig {
	config := &WebConfig{
		Address: address,
	}

	return config
}

// WebController ... web controller structure
type WebController struct {
	httpServer *http.Server
	Config     *Config
	started    bool
}

// NewWebController ... create a new WebController
func NewWebController(config *Config) (IWebController, error) {

	echo := echo.New()

	return &WebController{
		httpServer: echo,
		Config:     config,
	}, nil
}

// AddRoute ... adds a new route handler
func (manager *WebController) AddRoute(method string, route string, handler func(context echo.Context) error) {
	log.Infof("web, adding '%s' method to route '%s'", method, route)
	switch method {
	case "POST":
		manager.httpServer.POST(route, handler)
	case "PUT":
		manager.httpServer.PUT(route, handler)
	case "GET":
		manager.httpServer.GET(route, handler)
	case "HEAD":
		manager.httpServer.HEAD(route, handler)
	case "PATCH":
		manager.httpServer.PATCH(route, handler)
	case "OPTIONS":
		manager.httpServer.OPTIONS(route, handler)
	}
}

// Start ... starts the server
func (manager *WebController) Start() error {
	manager.started = true

	log.Infof("web, starting webserver [address:%s]", manager.Config.Address)
	if err := manager.httpServer.Start(manager.Config.Address); err != nil {
		log.Errorf("web, error starting webserver [address:%s], %s", manager.Config.Address, err.Error())
		manager.started = false
		return err
	}

	return nil
}

// Stop ... stops the server
func (manager *WebController) Stop() error {
	if !manager.started {
		manager.started = false
		return nil
	}

	return manager.httpServer.Server.Shutdown(context.Background())
}
