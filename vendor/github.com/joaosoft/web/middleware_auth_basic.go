package web

import (
	"github.com/joaosoft/auth-types/basic"
)

func MiddlewareCheckAuthBasic(user, password string) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) error {

			authHeader := ctx.Request.GetHeader(HeaderAuthorization)

			ok, err := basic.Check(authHeader, func(username string) (*basic.Credentials, error) {
				return &basic.Credentials{
					UserName: user,
					Password: password,
				}, nil
			})

			if !ok || err != nil {
				return ErrorInvalidAuthorization
			}

			return next(ctx)
		}
	}
}
