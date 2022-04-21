package web

import (
	"io"
)

func (w *Server) NewResponse(request *Request) *Response {
	return &Response{
		Base:                request.Base,
		FormData:            make(map[string]*FormData),
		MultiAttachmentMode: w.multiAttachmentMode,
		Boundary:            RandomBoundary(),
		Writer:              request.conn.(io.Writer),
		Status:              StatusNoContent,
		StatusText:          StatusText(StatusNoContent),
	}
}
