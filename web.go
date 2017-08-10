package mgr

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/joaosoft/go-Manager/web"
)

// -------------- WEB SERVERS --------------
// NewWEBServer ... creates a new web server
func (instance *Manager) NewWEBServer(config *web.Config) (web.IWebController, error) {
	log.Infof(fmt.Sprintf("web, creating web server"))
	return web.NewWebController(config)
}

// -------------- METHODS --------------
