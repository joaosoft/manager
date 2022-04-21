package web

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"strings"
	"time"
)

func (w *Server) NewRequest(conn net.Conn, server *Server) (*Request, error) {

	request := &Request{
		Base: Base{
			Server:      server,
			IP:          conn.RemoteAddr().String(),
			Headers:     make(Headers),
			Cookies:     make(Cookies),
			Params:      make(Params),
			UrlParams:   make(UrlParams),
			ContentType: ContentTypeEmpty,
			Charset:     CharsetUTF8,
			conn:        conn,
		},
		FormData:    make(map[string]*FormData),
		Attachments: make(map[string]*Attachment),
		Reader:      conn.(io.Reader),
	}

	return request, request.read()
}

func (r *Request) Bind(obj interface{}) error {
	contentType := r.GetContentType()

	if len(r.Body) == 0 || contentType == nil {
		return nil
	}

	switch *contentType {
	case ContentTypeApplicationJSON:
		if err := json.Unmarshal(r.Body, obj); err != nil {
			return err
		}
	case ContentTypeApplicationXML:
		if err := xml.Unmarshal(r.Body, obj); err != nil {
			return err
		}
	default:
		tmp := string(r.Body)
		obj = &tmp
	}
	return nil
}

func (r *Request) BindFormData(obj interface{}) error {
	if len(r.FormData) == 0 {
		return nil
	}

	data := make(map[string]string)
	for _, item := range r.FormData {
		if item.IsAttachment {
			continue
		}

		data[item.Name] = string(item.Body)
	}

	return readData(reflect.ValueOf(obj), data)
}

func (r *Request) read() error {
	reader := bufio.NewReader(r.conn)

	// header
	if err := r.readHeader(reader); err != nil {
		return err
	}

	// headers
	if err := r.readHeaders(reader); err != nil {
		return err
	}

	// boundary
	if r.Boundary != "" {
		r.handleBoundary(reader)
	} else {
		// body
		switch r.ContentType {
		case ContentTypeApplicationForm:
			if err := r.readUrlForm(reader); err != nil {
				return err
			}
		default:
			if err := r.readBody(reader); err != nil {
				return err
			}
		}

	}

	return nil
}

func (r *Request) readHeader(reader *bufio.Reader) error {

	// read one line (ended with \n or \r\n)
	r.conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	line, _, err := reader.ReadLine()
	if err != nil {
		return fmt.Errorf("invalid http request: %s", err)
	}

	if firstLine := bytes.SplitN(line, []byte(` `), 3); len(firstLine) < 3 {
		return errors.New("invalid http request")
	} else {
		r.Method = Method(firstLine[0])
		r.Address = NewAddress(string(firstLine[1]))
		r.Protocol = Protocol(firstLine[2])
		r.Params = r.Address.Params
	}

	return nil
}

func (r *Request) readHeaders(reader *bufio.Reader) error {
	for {
		r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
		line, _, err := reader.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}

		if split := bytes.SplitN(line, []byte(`: `), 2); len(split) > 0 {
			switch string(bytes.TrimSpace(bytes.Title(split[0]))) {
			case "Cookie":
				var cookieValue string
				splitCookie := bytes.Split(split[1], []byte(`=`))
				if len(splitCookie) > 1 {
					cookieValue = string(splitCookie[1])
				}
				r.Cookies[strings.Title(string(split[0]))] = Cookie{Name: string(splitCookie[0]), Value: cookieValue}
			case "Content-Type":
				if args := bytes.Split(split[1], []byte(`;`)); len(args) > 0 {
					split[1] = bytes.TrimSpace(args[0])
					for _, arg := range args {
						parm := bytes.Split(arg, []byte(`=`))
						switch string(bytes.TrimSpace(parm[0])) {
						case "boundary":
							r.Boundary = string(bytes.Replace(parm[1], []byte(`"`), []byte(``), -1))
						case "charset":
							r.Charset = Charset(bytes.Replace(parm[1], []byte(`"`), []byte(``), -1))
						}
					}
				}
				fallthrough
			default:
				r.Headers[strings.Title(string(split[0]))] = []string{string(bytes.TrimSpace(split[1]))}
			}
		}
	}

	return nil
}

func (r *Request) handleBoundary(reader *bufio.Reader) error {
	// read next line
	r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
	line, _, err := reader.ReadLine()
	if err != nil {
		return err
	}

	for {
		data := &Data{}
		formDataBody := bytes.NewBuffer(nil)

		for {
			content := bytes.SplitN(line, []byte(`: `), 2)
			switch string(bytes.Title(bytes.TrimSpace(content[0]))) {
			case "Content-Type":
				bytes.Split(content[1], []byte(`;`))
				data.ContentType = ContentType(content[1])

			case "Content-Disposition":
				split := bytes.Split(content[1], []byte(`;`))
				data.ContentDisposition = ContentDisposition(string(split[0]))
				for i := 1; i < len(split); i++ {
					parms := bytes.Split(split[i], []byte(`=`))
					switch string(bytes.TrimSpace(parms[0])) {
					case "name":
						if len(parms) > 1 {
							data.Name = string(bytes.Replace(parms[1], []byte(`"`), []byte(""), 2))
						}
					case "filename":
						if len(parms) > 1 {
							data.FileName = string(bytes.Replace(parms[1], []byte(`"`), []byte(""), 2))
						}
						data.IsAttachment = true
					}
				}
			}

			// read next line
			r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
			line, _, err = reader.ReadLine()
			if err != nil || len(line) == 0 {
				break
			}
		}

		for {
			r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
			line, err = reader.ReadSlice('\n')

			if !bytes.HasPrefix(line, []byte(fmt.Sprintf("--%s--", r.Boundary))) { // here we dont have a new line

				if err != nil {
					data.Body = formDataBody.Bytes()
					return err
				}

				if !data.IsAttachment {
					if line[len(line)-1] == '\n' {
						drop := 1
						if len(line) > 1 && line[len(line)-2] == '\r' {
							drop = 2
						}
						line = line[:len(line)-drop]
					}
				}
			}

			// is another boundary ?
			if bytes.HasPrefix(line, []byte(fmt.Sprintf("--%s", r.Boundary))) ||
				bytes.HasPrefix(line, []byte(fmt.Sprintf("--%s--", r.Boundary))) {
				// save formData
				data.Body = formDataBody.Bytes()
				key := data.Name
				if key == "" {
					key = data.FileName
				}

				if data.IsAttachment {
					r.Attachments[key] = &Attachment{
						Data: data,
					}
				} else {
					r.FormData[key] = &FormData{
						Data: data,
					}
				}

				// next formData
				data = &Data{}
				formDataBody = bytes.NewBuffer(nil)

				break
			} else {
				formDataBody.Write(line)
			}
		}

		if bytes.HasPrefix(line, []byte(fmt.Sprintf("--%s--", r.Boundary))) {
			return nil
		}
	}

	return nil
}

func (r *Request) readBody(reader *bufio.Reader) error {
	var buf bytes.Buffer
	for {
		r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		buf.Write(line)
	}
	r.Body = buf.Bytes()

	return nil
}

func (r *Request) readUrlForm(reader *bufio.Reader) error {
	var buf bytes.Buffer
	for {
		r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		buf.Write(line)
	}

	variables := bytes.Split(buf.Bytes(), []byte("&"))
	for _, variable := range variables {

		keyValues := bytes.Split(variable, []byte("="))
		lenKeyValues := len(keyValues)

		for i := 0; i < lenKeyValues; i += 2 {

			name := string(keyValues[0])
			r.FormData[name] = &FormData{
				&Data{
					Name: name,
					Body: keyValues[1],
				},
			}
		}
	}

	return nil
}
