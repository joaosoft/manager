package web

import (
	"fmt"
	"regexp"
)

func (c *Client) NewRequest(method Method, url string, contentType ContentType, headers Headers) (*Request, error) {

	// validate url
	regx := regexp.MustCompile(RegexForURL)
	if !regx.MatchString(url) {
		return nil, fmt.Errorf("invalid url [%s]", url)
	}

	address := NewAddress(url)
	params := address.Params

	if headers == nil {
		headers = make(Headers)
	}

	if contentType == ContentTypeEmpty {
		if c, ok := headers[HeaderContentType]; ok {
			contentType = ContentType(c[0])
		}
	}

	return &Request{
		Base: Base{
			Client:      c,
			Protocol:    ProtocolHttp1p1,
			Method:      method,
			Address:     address,
			Headers:     headers,
			Cookies:     make(Cookies),
			Params:      params,
			UrlParams:   make(UrlParams),
			Charset:     CharsetUTF8,
			ContentType: contentType,
		},
		FormData:            make(map[string]*FormData),
		Attachments:         make(map[string]*Attachment),
		MultiAttachmentMode: c.multiAttachmentMode,
		Boundary:            RandomBoundary(),
	}, nil
}
