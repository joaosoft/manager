package gomanager

import (
	"io"
)

type IGatewayManager interface {
	AddDefaultHeader(key, value string)
	Request(method string, endpoint string, headers map[string]string, body io.Reader) (int, []byte, error)
}
