package web

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/joaosoft/color"
	"github.com/joaosoft/logger"
)

type Server struct {
	name                string
	config              *ServerConfig
	isLogExternal       bool
	logger              logger.ILogger
	routes              Routes
	filters             Filters
	middlewares         []MiddlewareFunc
	listener            net.Listener
	errorhandler        ErrorHandler
	multiAttachmentMode MultiAttachmentMode
	started             bool
}

func NewServer(options ...ServerOption) (*Server, error) {
	config, err := NewServerConfig()

	service := &Server{
		name:                "server",
		logger:              logger.NewLogDefault("server", logger.WarnLevel),
		routes:              make(Routes),
		filters:             make(Filters),
		middlewares:         make([]MiddlewareFunc, 0),
		multiAttachmentMode: MultiAttachmentModeZip,
		config:              &config.Server,
	}

	if service.isLogExternal {
		// set logger of internal processes
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else {
		level, _ := logger.ParseLevel(service.config.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.Reconfigure(options...)

	if config.Server.Address == "" {
		port, err := GetFreePort()
		if err != nil {
			return nil, err
		}
		config.Server.Address = fmt.Sprintf(":%d", port)
	}

	if err = service.AddRoute(MethodGet, "/favicon.ico", service.handlerFile); err != nil {
		return nil, err
	}

	service.errorhandler = service.DefaultErrorHandler

	return service, nil
}

func (w *Server) AddMiddlewares(middlewares ...MiddlewareFunc) {
	w.middlewares = append(w.middlewares, middlewares...)
}

func (w *Server) AddFilter(pattern string, position Position, middleware MiddlewareFunc, method Method, methods ...Method) {
	w.filters.AddFilter(pattern, position, middleware, method, methods...)
}

func (w *Server) AddRoute(method Method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	w.routes[method] = append(w.routes[method], Route{
		Method:      method,
		Path:        path,
		Regex:       ConvertPathToRegex(path),
		Handler:     handler,
		Middlewares: middleware,
		Name:        GetFunctionName(handler),
	})

	return nil
}

func (w *Server) AddRoutes(route ...*Route) error {
	for _, r := range route {
		if err := w.AddRoute(r.Method, r.Path, r.Handler, r.Middlewares...); err != nil {
			return err
		}
	}
	return nil
}

func (w *Server) AddNamespace(path string, middlewares ...MiddlewareFunc) *Namespace {
	return &Namespace{
		Path:        path,
		Middlewares: middlewares,
		WebServer:   w,
	}
}

func (n *Namespace) AddRoutes(route ...*Route) error {
	for _, r := range route {
		if err := n.WebServer.AddRoute(r.Method, fmt.Sprintf("%s%s", n.Path, r.Path), r.Handler, append(r.Middlewares, n.Middlewares...)...); err != nil {
			return err
		}
	}
	return nil
}

func (n *Namespace) AddRoute(method Method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	return n.WebServer.AddRoute(method, fmt.Sprintf("%s%s", path, path), handler, append(middleware, n.Middlewares...)...)
}

func (w *Server) SetErrorHandler(handler ErrorHandler) error {
	w.errorhandler = handler
	return nil
}

func (w *Server) handleConnection(conn net.Conn) (err error) {
	var ctx *Context
	var length int
	var handlerRoute HandlerFunc
	var skipRouterHandler bool
	startTime := time.Now()

	defer func() {
		conn.Close()
	}()

	// read response from connection
	request, err := w.NewRequest(conn, w)
	if err != nil {
		w.logger.Errorf("error getting request: [%s]", err)
		return err
	}

	if w.logger.IsDebugEnabled() {
		if request.Body != nil {
			w.logger.Infof("[REQUEST BODY] [%s]", string(request.Body))
		}
	}

	// create response for request
	response := w.NewResponse(request)

	// create context with request and response
	ctx = NewContext(startTime, request, response)
	var route *Route

	// when options method, validate request route
	if request.Method == MethodOptions {
		if _, ok := w.routes[MethodOptions]; !ok {
			var method Method
			if val, ok := request.Headers[HeaderAccessControlRequestMethod]; ok {
				method = Method(val[0])
			}
			route, err = w.GetRoute(method, request.Address.Url)
			if err != nil {
				w.logger.Errorf("error getting route: [%s]", err)
				goto done
			}

			if err == nil && route != nil {
				ctx.Response.Headers[HeaderAccessControlAllowMethods] = []string{string(method)}
				ctx.Response.Headers[HeaderAccessControlAllowHeaders] = []string{strings.Join([]string{
					string(HeaderContentType),
					string(HeaderAccessControlAllowHeaders),
					string(HeaderAuthorization),
					string(HeaderXRequestedWith),
				}, ", ")}
				goto done
			}

			skipRouterHandler = true
		}
	}

	// route of the Server
	route, err = w.GetRoute(request.Method, request.Address.Url)
	if err != nil {
		w.logger.Errorf("error getting route: [%s]", err)
		goto done
	}

	// get url parameters
	if err = w.LoadUrlParms(request, route); err != nil {
		w.logger.Errorf("error loading url parameters: [%s]", err)
		goto done
	}

	// execute before filters
	handlerRoute = emptyHandler
	if after, ok := w.filters[PositionAfter]; ok {

		filters, err := w.GetMatchedFilters(after, request.Method, request.Address.Url)
		if err != nil {
			w.logger.Errorf("error getting route: [%s]", err)
			goto done
		}

		length = len(filters)
		for i, _ := range filters {
			if filters[length-1-i] != nil {
				handlerRoute = filters[length-1-i].Middleware(handlerRoute)
			}
		}
	}

	// route handler
	if !skipRouterHandler {
		handlerRoute = func(next HandlerFunc) HandlerFunc {
			return func(ctx *Context) error {
				if err := route.Handler(ctx); err != nil {
					return err
				}

				return next(ctx)
			}

		}(handlerRoute)
	}

	// execute before filters
	if between, ok := w.filters[PositionBetween]; ok {

		filters, err := w.GetMatchedFilters(between, request.Method, request.Address.Url)
		if err != nil {
			w.logger.Errorf("error getting route: [%s]", err)
			goto done
		}

		length = len(filters)
		for i, _ := range filters {
			if filters[length-1-i] != nil {
				handlerRoute = filters[length-1-i].Middleware(handlerRoute)
			}
		}
	}

	// execute middlewares
	length = len(w.middlewares)
	for i, _ := range w.middlewares {
		if w.middlewares[length-1-i] != nil {
			handlerRoute = w.middlewares[length-1-i](handlerRoute)
		}
	}

	// middleware's of the specific route
	length = len(route.Middlewares)
	for i, _ := range route.Middlewares {
		if route.Middlewares[length-1-i] != nil {
			handlerRoute = route.Middlewares[length-1-i](handlerRoute)
		}
	}

	// execute before filters
	if before, ok := w.filters[PositionBefore]; ok {

		filters, err := w.GetMatchedFilters(before, request.Method, request.Address.Url)
		if err != nil {
			w.logger.Errorf("error getting route: [%s]", err)
			goto done
		}

		length = len(filters)
		for i, _ := range filters {
			if filters[length-1-i] != nil {
				handlerRoute = filters[length-1-i].Middleware(handlerRoute)
			}
		}
	}

	// run handlers with middleware's
	if err = handlerRoute(ctx); err != nil {
		w.logger.Errorf("error executing handler: [%s]", err)
		goto done
	}

done:
	if err != nil {
		if er := w.errorhandler(ctx, err); er != nil {
			w.logger.Errorf("error writing error: [error: %s] %s", err, er)
		}
	}

	// write response to connection
	if err = ctx.Response.write(); err != nil {
		w.logger.Errorf("error writing response: [%s]", err)
	}

	fmt.Println(color.WithColor("Server[%s] Status[%d] Address[%s] Method[%s] Url[%s] Protocol[%s] Start[%s] Elapsed[%s]", color.FormatBold, color.ForegroundCyan, color.BackgroundBlack, w.name, ctx.Response.Status, ctx.Request.IP, ctx.Request.Method, ctx.Request.Address.Url, ctx.Request.Protocol, startTime.Format(TimeFormat), time.Since(startTime)))

	return nil
}

func ConvertPathToRegex(path string) string {

	var re = regexp.MustCompile(`:[a-zA-Z0-9\-_.:]+`)
	regx := re.ReplaceAllString(path, `[a-zA-Z0-9-_.:]+`)
	regx = strings.Replace(regx, "*", "(.+)", -1)

	return fmt.Sprintf("^%s$", regx)
}

func (w *Server) GetRoute(method Method, url string) (*Route, error) {

	for _, route := range w.routes[method] {
		if regx, err := regexp.Compile(route.Regex); err != nil {
			return nil, err
		} else {
			if regx.MatchString(url) {
				return &route, nil
			}
		}
	}

	return nil, ErrorNotFound
}

func (w *Server) GetMatchedFilters(filters map[Method][]*Filter, method Method, url string) ([]*Filter, error) {
	matched := make([]*Filter, 0)

	for _, filter := range filters[method] {
		if regx, err := regexp.Compile(filter.Regex); err != nil {
			return nil, err
		} else {
			if regx.MatchString(url) {
				matched = append(matched, filter)
			}
		}
	}

	for _, filter := range filters[MethodAny] {
		if regx, err := regexp.Compile(filter.Regex); err != nil {
			return nil, err
		} else {
			if regx.MatchString(url) {
				matched = append(matched, filter)
			}
		}
	}

	return matched, nil
}

func (w *Server) LoadUrlParms(request *Request, route *Route) error {

	routeUrl := strings.Split(route.Path, "/")
	url := strings.Split(request.Address.Url, "/")

	for i, name := range routeUrl {
		if name != url[i] {
			request.UrlParams[name[1:]] = []string{url[i]}
		}
	}

	return nil
}

func emptyHandler(ctx *Context) error {
	return nil
}

func (w *Server) Start(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	w.logger.Debug("executing Start")
	var err error

	w.listener, err = net.Listen("tcp", w.config.Address)
	if err != nil {
		w.logger.Errorf("error connecting to %s: %s", w.config.Address, err)
		return err
	}

	if w.config.Address == ":0" {
		split := strings.Split(w.listener.Addr().String(), ":")
		w.config.Address = fmt.Sprintf(":%s", split[len(split)-1])
	}

	fmt.Println(color.WithColor("http Server [%s] started on [%s]", color.FormatBold, color.ForegroundRed, color.BackgroundBlack, w.name, w.config.Address))

	w.started = true
	wg.Done()

	for {
		conn, err := w.listener.Accept()
		w.logger.Info("accepted connection")
		if err != nil {
			w.logger.Errorf("error accepting connection: %s", err)
			continue
		}

		if conn == nil {
			w.logger.Error("the connection isn't initialized")
			continue
		}

		go w.handleConnection(conn)
	}

	return err
}

func (w *Server) Started() bool {
	return w.started
}

func (w *Server) Stop(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	w.logger.Debug("executing Stop")

	if w.listener != nil {
		w.listener.Close()
	}

	w.started = false

	return nil
}
