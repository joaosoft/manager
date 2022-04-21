package web

import (
	"strings"
)

func MiddlewareOptions() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) error {

			if ctx.Request.Method == MethodOptions {
				var method Method
				if val, ok := ctx.Request.Headers[HeaderAccessControlRequestMethod]; ok {
					method = Method(val[0])
				} else {
					return ctx.Response.NoContent(StatusBadRequest)
				}

				route, err := ctx.Request.Server.GetRoute(method, ctx.Request.Address.Url)
				if err == nil && route != nil {
					ctx.Response.Headers[HeaderAccessControlAllowMethods] = []string{string(ctx.Request.Method)}
					ctx.Response.Headers[HeaderAccessControlAllowHeaders] = []string{strings.Join([]string{
						string(HeaderContentType),
						string(HeaderAccessControlAllowHeaders),
						string(HeaderAuthorization),
						string(HeaderXRequestedWith),
					}, ", ")}
				} else if err != ErrorNotFound {
					return ctx.Response.NoContent(StatusNotFound)
				} else {
					return ctx.Response.NoContent(StatusBadRequest)
				}
			}

			return next(ctx)
		}
	}
}
