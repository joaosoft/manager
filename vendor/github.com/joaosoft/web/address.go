package web

import (
	"fmt"
	"strings"
)

type Params map[string][]string

type Address struct {
	Full      string
	Schema    Schema
	Url       string
	ParamsUrl string
	Host      string
	Params    Params
}

func (p Params) String() string {
	var url string

	lenP := len(p)
	if lenP > 0 {
		url += "?"
	}

	i := 1
	for key, value := range p {
		url += fmt.Sprintf("%s=%s", key, value[0])

		if i < lenP {
			url += "&"
				i++
		}
	}

	return url
}

func NewAddress(url string) *Address {
	var tmp, full, schema, host, paramsUrl string
	var params = make(Params)

	tmp = url
	full = tmp // full

	split := strings.SplitN(tmp, "://", 2)
	if len(split) == 2 {
		schema = split[0] // schema
		tmp = split[1]
	}

	split = strings.SplitN(tmp, "/", 2)
	host = split[0] // host

	if len(split) == 2 {
		tmp = split[1]
		url = fmt.Sprintf("/%s", tmp) // url
	}

	// load query parameters
	paramsUrl = fmt.Sprintf("/%s", tmp) // params url
	if split := strings.SplitN(tmp, "?", 2); len(split) > 1 {
		url = fmt.Sprintf("/%s", split[0]) // url
		if parms := strings.Split(split[1], "&"); len(parms) > 0 {
			for _, parm := range parms {
				if p := strings.Split(parm, "="); len(p) > 1 {
					if split := strings.SplitN(p[1], ",", -1); len(split) > 0 {
						for _, s := range split {
							params[p[0]] = append(params[p[0]], s)
						}
					}
					params[p[0]] = append(params[p[0]], p[1])
				}
			}
		}
	}

	return &Address{
		Full:      full,
		Schema:    Schema(schema),
		Host:      host,
		Url:       url,
		ParamsUrl: paramsUrl,
		Params:    params,
	}
}
