package web

import (
	"fmt"
)

func (ctx *Context) Redirect(host string) error {
	if ctx.Request.Client == nil {
		client, err := NewClient(WithClientLogger(ctx.Request.Server.logger))
		if err != nil {
			return err
		}

		ctx.Request.Client = client
	}

	url := fmt.Sprintf("%s%s", ctx.Request.Address.Url, ctx.Request.Params)
	ctx.Request.Address = NewAddress(fmt.Sprintf("%s%s", host, url))

	response, err := ctx.Request.Send()
	if err != nil {
		return err
	}

	return ctx.Response.Bytes(response.Status, response.ContentType, response.Body)
}
