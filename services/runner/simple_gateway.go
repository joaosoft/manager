package gomanager

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/gommon/log"
	"google.golang.org/genproto/googleapis/spanner/admin/manager/v1"
)

// Config ...
type Config struct {
	Host string `json:"host"`
}

// NewConfig ... created a new gateway config
func NewConfig(host string) *Config {
	config := &Config{
		Host: host,
	}

	return config
}

// Gateway ...
type Gateway struct {
	config         *Config
	client         *http.Client
	defaultHeaders map[string]string
}

// NewGateway ...
func NewGateway(config *Config) *Gateway {
	gateway := &Gateway{
		config:         config,
		client:         &http.Client{},
		defaultHeaders: make(map[string]string),
	}

	return gateway
}

// AddDefaultHeader ...
func (gateway *Gateway) AddDefaultHeader(key, value string) {
	gateway.defaultHeaders[key] = value
	teste := http.Client{}
	teste.Do()
}

// Request ...
func (gateway *Gateway) Request(method string, endpoint string, headers map[string]string, body io.Reader) (int, []byte, error) {
	url := fmt.Sprintf("%s/%s", gateway.config.Host, endpoint)
	log.Info(fmt.Sprintf("gateway, url: %s", url))
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, nil, fmt.Errorf(fmt.Sprintf("gateway, error creating request [url:%s]", err.Error()), err)
	}

	for key, value := range gateway.defaultHeaders {
		req.Header.Add(key, value)
	}
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	response, err := gateway.client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf(fmt.Sprintf("gateway, error running request [url:%s]", err.Error()), err)
	}
	defer response.Body.Close()

	output, err := ioutil.ReadAll(response.Body)

	return response.StatusCode, output, err
}
