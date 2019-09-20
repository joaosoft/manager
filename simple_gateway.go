package manager

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/web"

	"fmt"
)

// Headers ...
type Headers map[string][]string

// SimpleGateway ...
type SimpleGateway struct {
	client *web.Client
	logger logger.ILogger
}

// NewSimpleGateway ...
func (manager *Manager) NewSimpleGateway() (IGateway, error) {
	client, err := web.NewClient()
	if err != nil {
		return nil, err
	}

	return &SimpleGateway{
		client: client,
		logger: manager.logger,
	}, nil
}

// Request ...
func (gateway *SimpleGateway) Request(method, host, endpoint string, contentType string, headers map[string][]string, body []byte) (int, []byte, error) {
	url := fmt.Sprintf("%s%s", host, endpoint)

	request, err := gateway.client.NewRequest(web.Method(method), url, web.ContentType(contentType), headers)
	if err != nil {
		panic(err)
	}

	if body != nil {
		request.WithBody(body)
	}

	response, err := request.Send()
	if err != nil {
		return 0, nil, err
	}

	return int(response.Status), response.Body, nil
}
