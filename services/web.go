package mgr

import (
	"fmt"

	"github.com/joaosoft/go-manager/services/web"
	"github.com/labstack/gommon/log"
)

// -------------- WEB SERVERS --------------
// NewWEBServer ... creates a new web server
func (instance *Manager) NewWEBServer(config *web.Config) (web.IWebController, error) {
	log.Infof(fmt.Sprintf("web, creating web server"))
	return web.NewWebController(config)
}

// -------------- METHODS --------------
