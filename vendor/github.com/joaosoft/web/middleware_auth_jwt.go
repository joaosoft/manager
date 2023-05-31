package web

import (
	"github.com/joaosoft/auth-types/jwt"
	"strings"

)

func MiddlewareCheckAuthJwt(keyFunc jwt.KeyFunc, checkFunc jwt.CheckFunc) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) error {
			authHeader := ctx.Request.GetHeader(HeaderAuthorization)
			token := strings.Replace(authHeader, "Bearer ", "", 1)

			ok, err := jwt.Check(token, keyFunc, checkFunc, jwt.Claims{}, false)

			if !ok || err != nil {
				return ErrorInvalidAuthorization
			}

			return next(ctx)
		}
	}
}
