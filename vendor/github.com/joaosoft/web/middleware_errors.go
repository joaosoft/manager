package web

import (
	"net/http"

	"github.com/joaosoft/errors"
)

var (
	ErrorInvalidAuthorization = errors.New(errors.LevelError, http.StatusUnauthorized, "invalid authorization")
)
