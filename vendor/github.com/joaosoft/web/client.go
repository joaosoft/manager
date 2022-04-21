package web

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/joaosoft/color"
	"github.com/joaosoft/logger"
)

type Client struct {
	config              *ClientConfig
	isLogExternal       bool
	logger              logger.ILogger
	dialer              net.Dialer
	multiAttachmentMode MultiAttachmentMode
}

func NewClient(options ...ClientOption) (*Client, error) {
	config, err := NewClientConfig()

	service := &Client{
		logger:              logger.NewLogDefault("client", logger.WarnLevel),
		multiAttachmentMode: MultiAttachmentModeZip,
		config:              &config.Client,
	}

	if service.isLogExternal {
		// set logger of internal processes
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else {
		level, _ := logger.ParseLevel(service.config.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	// create a new dialer to create connections
	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	service.dialer = dialer

	service.Reconfigure(options...)

	return service, nil
}
func (r *Request) Send() (*Response, error) {
	return r.Client.Send(r)
}

func (c *Client) Send(request *Request) (*Response, error) {
	startTime := time.Now()

	if c.logger.IsDebugEnabled() {
		if request.Body != nil {
			c.logger.Infof("[REQUEST BODY] [%s]", string(request.Body))
		}
	}

	c.logger.Debugf("executing [%s] request to [%s]", request.Method, request.Address.Full)

	address := request.Address.Host
	if request.Address.Schema != SchemaNone {
		address += fmt.Sprintf(":%s", request.Address.Schema)
	}

	var conn net.Conn
	var err error

	switch request.Address.Schema {
	case SchemaHttps:
		conn, err = tls.Dial("tcp", address, nil)
	default:
		conn, err = c.dialer.Dial("tcp", address)
	}

	if err != nil {
		return nil, err
	}

	body, err := request.build()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.Write(body)

	response, err := c.NewResponse(request.Method, request.Address, conn)

	if c.logger.IsDebugEnabled() {
		if response.Body != nil {
			c.logger.Infof("[RESPONSE BODY] [%s]", string(response.Body))
		}
	}

	fmt.Println(color.WithColor("Status[%d] Method[%s] Url[%s] on Start[%s] Elapsed[%s]", color.FormatBold, color.ForegroundCyan, color.BackgroundBlack, response.Status, request.Method, request.Address.Url, startTime.Format(TimeFormat), time.Since(startTime)))

	return response, err
}
