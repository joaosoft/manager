package web

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/flate"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (r *Response) Set(status Status, contentType ContentType, b []byte) error {
	r.Status = status
	r.ContentType = contentType
	r.Body = b
	return nil
}

func (r *Response) HTML(status Status, body string) error {
	r.Status = status
	r.SetContentType(ContentTypeTextHTML)
	r.Body = []byte(body)
	return nil
}

func (r *Response) Bytes(status Status, contentType ContentType, b []byte) error {
	r.Status = status
	r.SetContentType(contentType)
	r.Body = b
	return nil
}

func (r *Response) String(status Status, s string) error {
	r.Status = status
	r.SetContentType(ContentTypeTextPlain)
	r.Body = []byte(s)
	return nil
}

func (r *Response) JSON(status Status, i interface{}) error {
	var pretty bool
	if value, ok := r.UrlParams["pretty"]; ok {
		pretty, _ = strconv.ParseBool(value[0])
	}

	if pretty {
		return r.JSONPretty(status, i, "  ")
	}

	if b, err := json.Marshal(i); err != nil {
		return err
	} else {
		r.Status = status
		r.SetContentType(ContentTypeApplicationJSON)
		r.Body = b
	}

	return nil
}

func (r *Response) JSONPretty(status Status, i interface{}, indent string) error {
	if b, err := json.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.Status = status
		r.SetContentType(ContentTypeApplicationJSON)
		r.Body = b
	}
	return nil
}

func (r *Response) XML(status Status, i interface{}) error {
	var pretty bool
	if value, ok := r.UrlParams["pretty"]; ok {
		pretty, _ = strconv.ParseBool(value[0])
	}

	if pretty {
		return r.XMLPretty(status, i, "  ")
	}

	if b, err := xml.Marshal(i); err != nil {
		return err
	} else {
		r.Status = status
		r.SetContentType(ContentTypeApplicationXML)
		r.Body = b
	}
	return nil
}

func (r *Response) XMLPretty(status Status, i interface{}, indent string) error {
	if b, err := xml.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.Status = status
		r.SetContentType(ContentTypeApplicationXML)
		r.Body = b
	}
	return nil
}

func (r *Response) Stream(status Status, contentType ContentType, reader io.Reader) error {
	r.Status = status
	r.SetContentType(contentType)
	if _, err := io.Copy(r.Writer, reader); err != nil {
		return err
	}
	return nil
}

func (r *Response) File(status Status, name string, body []byte) error {
	r.Status = status
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.SetContentType(contentType)
	r.SetCharset(charset)
	r.Body = body
	return nil
}

func (r *Response) Attachment(name string, body []byte) error {
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.FormData[name] = &FormData{
		Data: &Data{
			ContentDisposition: ContentDispositionAttachment,
			ContentType:        contentType,
			Charset:            charset,
			FileName:           name,
			Name:               name,
			Body:               body,
		},
	}
	return nil
}

func (r *Response) Inline(name string, body []byte) error {
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.FormData[name] = &FormData{
		Data: &Data{
			ContentDisposition: ContentDispositionInline,
			ContentType:        contentType,
			Charset:            charset,
			FileName:           name,
			Name:               name,
			Body:               body,
		},
	}
	return nil
}

func (r *Response) NoContent(status Status) error {
	r.Status = status
	return nil
}

func (r *Response) SetFormData(name string, value string) {
	r.FormData[name] = &FormData{
		Data: &Data{
			ContentDisposition: ContentDispositionFormData,
			ContentType:        ContentTypeTextPlain,
			Charset:            CharsetUTF8,
			Name:               name,
			Body:               []byte(value),
			IsAttachment:       false,
		},
	}
}

func (r *Response) GetFormDataBytes(name string) []byte {
	if value, ok := r.FormData[name]; ok {
		return value.Body
	}

	return nil
}

func (r *Response) GetFormDataString(name string) string {
	if value, ok := r.FormData[name]; ok {
		return string(value.Body)
	}

	return ""
}

func (r *Response) Bind(i interface{}) error {
	contentType := r.GetContentType()

	if len(r.Body) == 0 || contentType == nil {
		return nil
	}

	switch *contentType {
	case ContentTypeApplicationJSON:
		if err := json.Unmarshal(r.Body, i); err != nil {
			return err
		}
	case ContentTypeApplicationXML:
		if err := xml.Unmarshal(r.Body, i); err != nil {
			return err
		}
	default:
		tmp := string(r.Body)
		i = &tmp
	}
	return nil
}

func (r *Response) BindFormData(obj interface{}) error {
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

func setValue(kind reflect.Kind, obj reflect.Value, newValue string) error {

	if !obj.CanAddr() {
		return nil
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, _ := strconv.Atoi(newValue)
		obj.SetInt(int64(v))
	case reflect.Float32, reflect.Float64:
		v, _ := strconv.ParseFloat(newValue, 64)
		obj.SetFloat(v)
	case reflect.String:
		obj.SetString(newValue)
	case reflect.Bool:
		v, _ := strconv.ParseBool(newValue)
		obj.SetBool(v)
	}

	return nil
}

func (r *Response) read() error {
	reader := bufio.NewReader(r.Reader)

	// header
	if err := r.readHeader(reader); err != nil {
		return err
	}

	// headers
	if err := r.readHeaders(reader); err != nil {
		return err
	}

	// body
	if _, ok := MethodHasBody[r.Method]; ok {

		// boundary
		if r.Boundary != "" {
			r.readBoundary(reader)
		} else {

			// body
			if err := r.readBody(reader); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Response) readHeader(reader *bufio.Reader) error {

	// read one line (ended with \n or \r\n)
	r.conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	line, _, err := reader.ReadLine()
	if err != nil {
		return fmt.Errorf("invalid http send: %s", err)
	}

	if firstLine := bytes.SplitN(line, []byte(` `), 3); len(firstLine) < 3 {
		return errors.New("invalid http send")
	} else {
		status, err := strconv.Atoi(string(firstLine[1]))
		if err != nil {
			return fmt.Errorf("invalid http response [%s]", string(line))
		}

		r.Protocol = Protocol(firstLine[0])
		r.Status = Status(status)
		r.StatusText = string(firstLine[2])
	}

	return nil
}

func (r *Response) readHeaders(reader *bufio.Reader) error {
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
					r.ContentType = ContentType(split[1])
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
				r.Headers[strings.Title(string(split[0]))] = []string{string(split[1])}
			}
		}
	}

	return nil
}

func (r *Response) readBoundary(reader *bufio.Reader) error {

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
				contentDisposition := bytes.Split(content[1], []byte(`;`))
				data.ContentDisposition = ContentDisposition(string(contentDisposition[0]))
				for i := 1; i < len(contentDisposition); i++ {
					parms := bytes.Split(contentDisposition[i], []byte(`=`))
					switch string(bytes.TrimSpace(parms[0])) {
					case "name":
						data.Name = string(bytes.Replace(parms[1], []byte(`"`), []byte(""), 2))
					case "filename":
						data.FileName = string(bytes.Replace(parms[1], []byte(`"`), []byte(""), 2))
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
				// save data
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

				// next data
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

func (r *Response) readBody(reader *bufio.Reader) error {
	var buf bytes.Buffer
	var encoding = EncodingNone

	if enc, ok := r.Headers[HeaderTransferEncoding]; ok {
		encoding = Encoding(enc[0])
	}

	switch encoding {
	case EncodingChunked:
		var size uint64

		for {
			r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
			line, _, err := reader.ReadLine()
			if err != nil {
				break
			}

			size, _ = parseHexUint(line)
			if size == 0 {
				break
			}

			chunk := make([]byte, size)
			_, err = reader.Read(chunk)
			if err != nil {
				break
			}

			buf.Write(chunk)
		}
	default:
		for {
			r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
			line, _, err := reader.ReadLine()
			if err != nil {
				break
			}

			buf.Write(line)
		}
	}

	r.Body = buf.Bytes()

	return nil
}


func (r *Response) write() error {
	var buf bytes.Buffer
	var lenFormData = len(r.FormData)

	if headers, err := r.buildHeaders(); err != nil {
		return err
	} else {
		buf.Write(headers)
	}

	if lenFormData > 0 {
		switch r.MultiAttachmentMode {
		case MultiAttachmentModeBoundary:
			if body, err := r.buildBody(); err != nil {
				return err
			} else {
				buf.Write(body)
			}
			if body, err := r.buildBoundaryAttachments(); err != nil {
				return err
			} else {
				buf.Write(body)
			}
		case MultiAttachmentModeZip:
			if lenFormData > 1 {
				if body, err := r.buildZippedAttachments(); err != nil {
					return err
				} else {
					buf.Write(body)
				}
			} else {
				if body, err := r.buildSingleAttachment(); err != nil {
					return err
				} else {
					buf.Write(body)
				}
			}
		}
	} else {
		if body, err := r.buildBody(); err != nil {
			return err
		} else {
			buf.Write(body)

			if r.Server.logger.IsDebugEnabled() {
				r.Server.logger.Infof("[RESPONSE BODY] [%s]", string(body))
			}
		}
	}

	r.conn.Write(buf.Bytes())

	return nil
}

func (r *Response) buildHeaders() ([]byte, error) {
	var buf bytes.Buffer
	lenFormData := len(r.FormData)

	r.Headers[HeaderServer] = []string{"Server"}
	r.Headers[HeaderDate] = []string{time.Now().Format(HeaderTimeFormat)}
	r.Headers[HeaderAccessControlAllowCredentials] = []string{"true"}
	if val, ok := r.Headers[HeaderOrigin]; ok {
		r.Headers[HeaderAccessControlAllowOrigin] = val
	}

	// header
	buf.WriteString(fmt.Sprintf("%s %d %s\r\n", r.Protocol, r.Status, StatusText(r.Status)))

	if lenFormData > 0 {

		switch r.MultiAttachmentMode {
		case MultiAttachmentModeBoundary:
			r.Headers[HeaderContentType] = []string{fmt.Sprintf("%s; boundary=%s; charset=%s", ContentTypeMultipartFormData, r.Boundary, r.Charset)}
		case MultiAttachmentModeZip:
			var name = "attachments"
			var fileName = "attachments.zip"
			var contentType = ContentTypeApplicationZip
			var charset = r.Charset

			if lenFormData == 1 {
				for _, attachment := range r.FormData {
					name = attachment.Name
					fileName = attachment.FileName
					contentType = attachment.ContentType
					if attachment.Charset != "" {
						charset = attachment.Charset
					}
					break
				}
			}
			r.Headers[HeaderContentType] = []string{fmt.Sprintf("%s; attachment; name=%q; filename=%q; charset=%s", contentType, name, fileName, charset)}
		}
	} else {
		r.Headers[HeaderContentType] = []string{string(r.ContentType)}
		r.Headers[HeaderContentLength] = []string{strconv.Itoa(len(r.Body))}
	}

	for key, value := range r.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value[0]))
	}

	buf.WriteString("\r\n")

	return buf.Bytes(), nil
}

func (r *Response) buildBody() ([]byte, error) {
	var buf bytes.Buffer

	if MethodHasBody[r.Method] {
		buf.Write(r.Body)
		if r.MultiAttachmentMode == MultiAttachmentModeBoundary && len(r.FormData) > 0 {
			buf.WriteString("\r\n\r\n")
		}
	}

	return buf.Bytes(), nil
}

func (r *Response) buildSingleAttachment() ([]byte, error) {
	for _, formData := range r.FormData {
		return formData.Body, nil
	}
	return []byte{}, nil
}

func (r *Response) buildBoundaryAttachments() ([]byte, error) {
	var buf bytes.Buffer

	if len(r.FormData) == 0 {
		return buf.Bytes(), nil
	}

	lenF := len(r.FormData)
	i := 0

	for _, formData := range r.FormData {
		buf.WriteString(fmt.Sprintf("--%s\r\n", r.Boundary))
		if formData.IsAttachment {
			buf.WriteString(fmt.Sprintf("%s: %s; name=%q; filename=%q\r\n", HeaderContentDisposition, formData.ContentDisposition, formData.Name, formData.FileName))
		} else {
			buf.WriteString(fmt.Sprintf("%s: %s; name=%q\r\n", HeaderContentDisposition, formData.ContentDisposition, formData.Name))

		}
		buf.WriteString(fmt.Sprintf("%s: %s\r\n\r\n", HeaderContentType, formData.ContentType))
		buf.Write(formData.Body)

		if i < lenF-1 {
			buf.WriteString("\r\n")
		}
		i++
	}

	buf.WriteString(fmt.Sprintf("\r\n--%s--", r.Boundary))

	return buf.Bytes(), nil
}

func (r *Response) buildZippedAttachments() ([]byte, error) {
	// create a buffer to write our archive
	buf := new(bytes.Buffer)

	if len(r.FormData) == 0 {
		return buf.Bytes(), nil
	}

	// create a new zip archive
	w := zip.NewWriter(buf)

	// register a custom deflate compressor to override the default Deflate compressor with a higher compression level
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	for _, attachment := range r.FormData {
		f, err := w.Create(attachment.FileName)
		if err != nil {
			return buf.Bytes(), err
		}
		_, err = f.Write([]byte(attachment.Body))
		if err != nil {
			return buf.Bytes(), err
		}
	}

	err := w.Close()
	if err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}


