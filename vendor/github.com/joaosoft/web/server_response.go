package web

import (
	"io"
)

func (s *Server) NewResponse(request *Request) *Response {
	return &Response{
		Base:                request.Base,
		FormData:            make(map[string]*FormData),
		MultiAttachmentMode: s.multiAttachmentMode,
		Boundary:            RandomBoundary(),
		Writer:              request.conn.(io.Writer),
		Status:              StatusNoContent,
		StatusText:          StatusText(StatusNoContent),
	}
}
