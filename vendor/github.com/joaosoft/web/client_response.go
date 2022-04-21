package web

import (
	"io"
	"net"
)

func (c *Client) NewResponse(method Method, address *Address, conn net.Conn) (*Response, error) {

	response := &Response{
		Base: Base{
			Client:    c,
			Method:    method,
			Address:   address,
			Headers:   make(Headers),
			Cookies:   make(Cookies),
			Params:    make(Params),
			UrlParams: make(UrlParams),
			Charset:   CharsetUTF8,
			conn:      conn,
		},
		FormData:    make(map[string]*FormData),
		Attachments: make(map[string]*Attachment),
		Reader:      conn.(io.Reader),
	}

	return response, response.read()
}