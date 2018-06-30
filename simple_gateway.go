package manager

import (
	"io/ioutil"
	"net/http"

	"bytes"
	"fmt"
)

// Headers ...
type Headers map[string][]string

// SimpleGateway ...
type SimpleGateway struct {
	client *http.Client
}

// NewSimpleGateway ...
func NewSimpleGateway() IGateway {
	return &SimpleGateway{
		client: &http.Client{},
	}
}

// Request ...
func (gateway *SimpleGateway) Request(method, host, endpoint string, headers map[string][]string, body []byte) (int, []byte, error) {
	url := fmt.Sprintf("%s%s", host, endpoint)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return 0, nil, err
	}

	if headers != nil {
		for key, value := range headers {
			log.Infof("adding header with [ name: %s, value: %s ]", key, value)
			req.Header.Set(key, value[0])
		}
	}

	response, err := gateway.client.Do(req)

	var bodyResponse []byte

	if response != nil {
		defer response.Body.Close()
		bodyResponse, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return response.StatusCode, nil, err
		}
	}

	if err != nil {
		return 0, bodyResponse, err
	}

	return response.StatusCode, bodyResponse, nil
}
