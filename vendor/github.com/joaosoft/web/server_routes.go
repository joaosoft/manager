package web

type Routes map[Method][]Route

type Route struct {
	Method      Method
	Path        string
	Regex       string
	Name        string
	Handler     HandlerFunc
	Middlewares []MiddlewareFunc
}

func NewRoute(method Method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
	return &Route{
		Method:      method,
		Path:        path,
		Regex:       ConvertPathToRegex(path),
		Handler:     handler,
		Middlewares: middleware,
		Name:        GetFunctionName(handler),
	}
}
