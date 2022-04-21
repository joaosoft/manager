package web

import (
	"net/http"

	"github.com/joaosoft/errors"
)

var (
	ErrorInvalidAuthorization = errors.New(errors.ErrorLevel, http.StatusUnauthorized, "invalid authorization")
)
