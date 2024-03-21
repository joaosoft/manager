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
		logger:              logger.NewLogDefault("server", logger.LevelWarn),
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

func (s *Server) AddMiddlewares(middlewares ...MiddlewareFunc) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) AddFilter(pattern string, position Position, middleware MiddlewareFunc, method Method, methods ...Method) {
	s.filters.AddFilter(pattern, position, middleware, method, methods...)
}

func (s *Server) AddRoute(method Method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) error {
	s.routes[method] = append(s.routes[method], Route{
		Method:      method,
		Path:        path,
		Regex:       ConvertPathToRegex(path),
		Handler:     handler,
		Middlewares: middleware,
		Name:        GetFunctionName(handler),
	})

	return nil
}

func (s *Server) AddRoutes(route ...*Route) error {
	for _, r := range route {
		if err := s.AddRoute(r.Method, r.Path, r.Handler, r.Middlewares...); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) AddNamespace(path string, middlewares ...MiddlewareFunc) *Namespace {
	return &Namespace{
		Path:        path,
		Middlewares: middlewares,
		WebServer:   s,
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

func (s *Server) SetErrorHandler(handler ErrorHandler) error {
	s.errorhandler = handler
	return nil
}

func (s *Server) handleConnection(conn net.Conn) (err error) {
	var ctx *Context
	var length int
	var handlerRoute HandlerFunc
	var skipRouterHandler bool
	startTime := time.Now()

	defer func() {
		conn.Close()
	}()

	// read response from connection
	request, err := s.NewRequest(conn, s)
	if err != nil {
		s.logger.Errorf("error getting request: [%s]", err)
		return err
	}

	if s.logger.IsDebugEnabled() {
		if request.Body != nil {
			s.logger.Infof("[REQUEST BODY] [%s]", string(request.Body))
		}
	}

	// create response for request
	response := s.NewResponse(request)

	// create context with request and response
	ctx = NewContext(startTime, request, response)
	var route *Route

	// when options method, validate request route
	if request.Method == MethodOptions {
		if _, ok := s.routes[MethodOptions]; !ok {
			var method Method
			if val, ok := request.Headers[HeaderAccessControlRequestMethod]; ok {
				method = Method(val[0])
			}
			route, err = s.GetRoute(method, request.Address.Url)
			if err != nil {
				s.logger.Errorf("error getting route: [%s]", err)
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
	route, err = s.GetRoute(request.Method, request.Address.Url)
	if err != nil {
		s.logger.Errorf("error getting route: [%s]", err)
		goto done
	}

	// get url parameters
	if err = s.LoadUrlParms(request, route); err != nil {
		s.logger.Errorf("error loading url parameters: [%s]", err)
		goto done
	}

	// execute before filters
	handlerRoute = emptyHandler
	if after, ok := s.filters[PositionAfter]; ok {

		filters, err := s.GetMatchedFilters(after, request.Method, request.Address.Url)
		if err != nil {
			s.logger.Errorf("error getting route: [%s]", err)
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
	if between, ok := s.filters[PositionBetween]; ok {

		filters, err := s.GetMatchedFilters(between, request.Method, request.Address.Url)
		if err != nil {
			s.logger.Errorf("error getting route: [%s]", err)
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
	length = len(s.middlewares)
	for i, _ := range s.middlewares {
		if s.middlewares[length-1-i] != nil {
			handlerRoute = s.middlewares[length-1-i](handlerRoute)
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
	if before, ok := s.filters[PositionBefore]; ok {

		filters, err := s.GetMatchedFilters(before, request.Method, request.Address.Url)
		if err != nil {
			s.logger.Errorf("error getting route: [%s]", err)
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
		s.logger.Errorf("error executing handler: [%s]", err)
		goto done
	}

done:
	if err != nil {
		if er := s.errorhandler(ctx, err); er != nil {
			s.logger.Errorf("error writing error: [error: %s] %s", err, er)
		}
	}

	// write response to connection
	if err = ctx.Response.write(); err != nil {
		s.logger.Errorf("error writing response: [%s]", err)
	}

	fmt.Println(color.WithColor("Server[%s] Status[%d] Address[%s] Method[%s] Url[%s] Protocol[%s] Start[%s] Elapsed[%s]", color.FormatBold, color.ForegroundCyan, color.BackgroundBlack, s.name, ctx.Response.Status, ctx.Request.IP, ctx.Request.Method, ctx.Request.Address.Url, ctx.Request.Protocol, startTime.Format(TimeFormat), time.Since(startTime)))

	return nil
}

func ConvertPathToRegex(path string) string {

	var re = regexp.MustCompile(`:[a-zA-Z0-9\-_.:]+`)
	regx := re.ReplaceAllString(path, `[a-zA-Z0-9-_.:]+`)
	regx = strings.Replace(regx, "*", "(.+)", -1)

	return fmt.Sprintf("^%s$", regx)
}

func (s *Server) GetRoute(method Method, url string) (*Route, error) {

	for _, route := range s.routes[method] {
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

func (s *Server) GetMatchedFilters(filters map[Method][]*Filter, method Method, url string) ([]*Filter, error) {
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

func (s *Server) LoadUrlParms(request *Request, route *Route) error {

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

func (s *Server) Start(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	s.logger.Debug("executing Start")
	var err error

	s.listener, err = net.Listen("tcp", s.config.Address)
	if err != nil {
		s.logger.Errorf("error connecting to %s: %s", s.config.Address, err)
		return err
	}

	if s.config.Address == ":0" {
		split := strings.Split(s.listener.Addr().String(), ":")
		s.config.Address = fmt.Sprintf(":%s", split[len(split)-1])
	}

	fmt.Println(color.WithColor("http Server [%s] started on [%s]", color.FormatBold, color.ForegroundRed, color.BackgroundBlack, s.name, s.config.Address))

	s.started = true
	wg.Done()

	for {
		conn, err := s.listener.Accept()
		s.logger.Info("accepted connection")
		if err != nil {
			s.logger.Errorf("error accepting connection: %s", err)
			continue
		}

		if conn == nil {
			s.logger.Error("the connection isn't initialized")
			continue
		}

		go s.handleConnection(conn)
	}

	return err
}

func (s *Server) Started() bool {
	return s.started
}

func (s *Server) Stop(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	s.logger.Debug("executing Stop")

	if s.listener != nil {
		s.listener.Close()
	}

	s.started = false

	return nil
}

func (s *Server) Config() *ServerConfig {
	return s.config
}
