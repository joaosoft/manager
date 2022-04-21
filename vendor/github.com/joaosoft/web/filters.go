package web

type Filters map[Position]map[Method][]*Filter

type Filter struct {
	Method     Method
	Position   Position
	Pattern    string
	Regex      string
	Middleware MiddlewareFunc
}

func NewFilter(method Method, pattern string, position Position, middleware MiddlewareFunc) *Filter {
	return &Filter{
		Method:     method,
		Position:   position,
		Pattern:    pattern,
		Regex:      ConvertPathToRegex(pattern),
		Middleware: middleware,
	}
}

func (f Filters) AddFilter(pattern string, position Position, middleware MiddlewareFunc, method Method, methods ...Method) {
	if _, ok := f[position]; !ok {
		f[position] = make(map[Method][]*Filter)
	}

	f[position][method] = append(f[position][method], NewFilter(method, pattern, position, middleware))

	for _, m := range methods {
		f[position][m] = append(f[position][m], NewFilter(m, pattern, position, middleware))
	}
}
