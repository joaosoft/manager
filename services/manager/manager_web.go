package gomanager

import "github.com/labstack/echo"

// IWebController ... web controller interface
type IWebController interface {
	AddRoute(method string, route string, handler func(context echo.Context) error)
	Start() error
	Stop() error
}
