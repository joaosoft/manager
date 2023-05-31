package web

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joaosoft/auth-types/basic"
	"github.com/joaosoft/auth-types/jwt"
)

func (r *Request) Set(contentType ContentType, b []byte) error {
	r.ContentType = contentType
	r.Body = b
	return nil
}

func (r *Request) HTML(body string) error {
	r.SetContentType(ContentTypeTextHTML)
	r.Body = []byte(body)
	return nil
}

func (r *Request) Bytes(contentType ContentType, b []byte) error {
	r.SetContentType(contentType)
	r.Body = b
	return nil
}

func (r *Request) String(s string) error {
	r.SetContentType(ContentTypeTextPlain)
	r.Body = []byte(s)
	return nil
}

func (r *Request) JSON(i interface{}) error {
	var pretty bool
	if value, ok := r.UrlParams["pretty"]; ok {
		pretty, _ = strconv.ParseBool(value[0])
	}

	if pretty {
		return r.JSONPretty(i, "  ")
	}

	if b, err := json.Marshal(i); err != nil {
		return err
	} else {
		r.SetContentType(ContentTypeApplicationJSON)
		r.Body = b
	}

	return nil
}

func (r *Request) JSONPretty(i interface{}, indent string) error {
	if b, err := json.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(ContentTypeApplicationJSON)
		r.Body = b
	}
	return nil
}

func (r *Request) XML(i interface{}) error {
	var pretty bool
	if value, ok := r.UrlParams["pretty"]; ok {
		pretty, _ = strconv.ParseBool(value[0])
	}

	if pretty {
		return r.XMLPretty(i, "  ")
	}

	if b, err := xml.Marshal(i); err != nil {
		return err
	} else {
		r.SetContentType(ContentTypeApplicationXML)
		r.Body = b
	}
	return nil
}

func (r *Request) XMLPretty(i interface{}, indent string) error {
	if b, err := xml.MarshalIndent(i, "", indent); err != nil {
		return err
	} else {
		r.SetContentType(ContentTypeApplicationXML)
		r.Body = b
	}
	return nil
}

func (r *Request) Stream(contentType ContentType, reader io.Reader) error {
	r.SetContentType(contentType)
	if _, err := io.Copy(r.Writer, reader); err != nil {
		return err
	}
	return nil
}

func (r *Request) File(name string, body []byte) error {
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.SetContentType(contentType)
	r.SetCharset(charset)
	r.Body = body
	return nil
}

func (r *Request) Attachment(name string, body []byte) error {
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.FormData[name] = &FormData{
		Data: &Data{
			ContentDisposition: ContentDispositionAttachment,
			ContentType:        contentType,
			Charset:            charset,
			FileName:           name,
			Name:               name,
			Body:               body,
			IsAttachment:       true,
		},
	}
	return nil
}

func (r *Request) Inline(name string, body []byte) error {
	contentType, charset := DetectContentType(filepath.Ext(name), body)
	r.FormData[name] = &FormData{
		Data: &Data{
			ContentDisposition: ContentDispositionInline,
			ContentType:        contentType,
			Charset:            charset,
			FileName:           name,
			Name:               name,
			Body:               body,
			IsAttachment:       true,
		},
	}
	return nil
}

func (r *Request) SetFormData(name string, value string) {
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

func (r *Request) GetFormDataBytes(name string) []byte {
	if value, ok := r.FormData[name]; ok {
		return value.Body
	}

	return nil
}

func (r *Request) GetFormDataString(name string) string {
	if value, ok := r.FormData[name]; ok {
		return string(value.Body)
	}

	return ""
}

func (r *Request) WithBody(body []byte) *Request {
	r.Body = body

	return r
}

func (r *Request) WithAuthBasic(username, password string) (*Request, error) {
	r.SetHeader(HeaderAuthorization, []string{basic.Generate(username, password)})

	return r, nil
}

func (r *Request) WithAuthJwt(claims jwt.Claims, key interface{}) (*Request, error) {
	token, err := jwt.New(jwt.SignatureHS384).Generate(claims, key)
	if err != nil {
		return r, err
	}

	r.SetHeader(HeaderAuthorization, []string{token})

	return r, nil
}

func (r *Request) WithContentType(contentType ContentType) *Request {
	r.ContentType = contentType

	return r
}

func (r *Request) build() ([]byte, error) {
	var buf bytes.Buffer
	var lenAttachments = len(r.Attachments)

	if headers, err := r.buildHeaders(); err != nil {
		return nil, err
	} else {
		buf.Write(headers)
	}

	if lenAttachments > 0 {
		switch r.MultiAttachmentMode {
		case MultiAttachmentModeBoundary:
			if body, err := r.buildBody(); err != nil {
				return nil, err
			} else {
				buf.Write(body)
			}
			if body, err := r.buildBoundaries(); err != nil {
				return nil, err
			} else {
				buf.Write(body)
			}
		case MultiAttachmentModeZip:
			if lenAttachments > 1 {
				if body, err := r.buildZippedAttachments(); err != nil {
					return nil, err
				} else {
					buf.Write(body)
				}
			} else {
				if body, err := r.buildSingleAttachment(); err != nil {
					return nil, err
				} else {
					buf.Write(body)
				}
			}

			if body, err := r.buildBoundaries(); err != nil {
				return nil, err
			} else {
				buf.Write(body)
			}
		}
	} else {
		switch r.ContentType {
		case ContentTypeMultipartFormData, ContentTypeMultipartMixed:
			if body, err := r.buildBoundaries(); err != nil {
				return nil, err
			} else {
				buf.Write(body)
			}
		case ContentTypeApplicationForm:
			if urlForm, err := r.buildUrlForm(); err != nil {
				return nil, err
			} else {
				buf.Write(urlForm)
			}
		default:
			if body, err := r.buildBody(); err != nil {
				return nil, err
			} else {
				buf.Write(body)
			}
		}
	}

	return buf.Bytes(), nil
}

func (r *Request) buildHeaders() ([]byte, error) {
	var buf bytes.Buffer
	lenFormData := len(r.FormData)

	// header
	buf.WriteString(fmt.Sprintf("%s %s %s\r\n", r.Method, r.Address.ParamsUrl, r.Protocol))

	// headers
	r.Headers[HeaderHost] = []string{r.Address.Host}
	if _, ok := r.Headers[HeaderUserAgent]; !ok {
		r.Headers[HeaderUserAgent] = []string{"client"}
	}
	r.Headers[HeaderDate] = []string{time.Now().Format(HeaderTimeFormat)}

	if lenFormData > 0 {

		switch r.MultiAttachmentMode {
		case MultiAttachmentModeBoundary:
			r.Headers[HeaderContentType] = []string{fmt.Sprintf("%s; boundary=%s; charset=%s", ContentTypeMultipartFormData, r.Boundary, r.Charset)}
		case MultiAttachmentModeZip:
			if len(r.FormData) == 0 {
				var name = "attachments"
				var fileName = "attachments.zip"
				var contentType = ContentTypeApplicationZip
				var charset = r.Charset

				if lenFormData == 1 {
					for _, formData := range r.Attachments {
						name = formData.Name
						fileName = formData.FileName
						contentType = formData.ContentType
						if formData.Charset != "" {
							charset = formData.Charset
						}
						break
					}
				}
				r.Headers[HeaderContentType] = []string{fmt.Sprintf("%s; attachment; name=%q; filename=%q; charset=%s", contentType, name, fileName, charset)}
			} else {
				r.Headers[HeaderContentType] = []string{fmt.Sprintf("%s; boundary=%s; charset=%s", ContentTypeMultipartFormData, r.Boundary, r.Charset)}
			}
		}
	} else {
		if r.ContentType != ContentTypeEmpty {
			r.Headers[HeaderContentType] = []string{string(r.ContentType)}
		}
		lenBody := len(r.Body)
		if lenBody > 0 {
			r.Headers[HeaderContentLength] = []string{strconv.Itoa(lenBody)}
		}
	}

	for key, value := range r.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value[0]))
	}

	buf.WriteString("\r\n")

	return buf.Bytes(), nil
}

func (r *Request) buildBody() ([]byte, error) {
	var buf bytes.Buffer

	if MethodHasBody[r.Method] {
		buf.Write(r.Body)
		if r.MultiAttachmentMode == MultiAttachmentModeBoundary && len(r.FormData) > 0 {
			buf.WriteString("\r\n\r\n")
		}
	}

	return buf.Bytes(), nil
}

func (r *Request) buildUrlForm() ([]byte, error) {
	var buf bytes.Buffer

	lenI := len(r.FormData)
	i := 0
	for _, formData := range r.FormData {
		if formData.IsAttachment {
			continue
		}

		buf.WriteString(fmt.Sprintf("%s=%s", formData.Name, string(formData.Body)))

		if i < lenI-1 {
			buf.WriteString("&")
		}
		i++
	}

	return buf.Bytes(), nil
}

func (r *Request) buildSingleAttachment() ([]byte, error) {
	for _, attachment := range r.FormData {
		return attachment.Body, nil
	}
	return []byte{}, nil
}

func (r *Request) buildBoundaries() ([]byte, error) {
	var buf bytes.Buffer

	switch r.ContentType {
	case ContentTypeMultipartFormData, ContentTypeMultipartMixed:
		bufFormData, err := r.buildFormData()
		if err != nil {
			return buf.Bytes(), err
		}
		buf.Write(bufFormData)

	default:
		if r.MultiAttachmentMode != MultiAttachmentModeZip {
			bufAttachments, err := r.buildAttachments()
			if err != nil {
				return buf.Bytes(), err
			}
			buf.Write(bufAttachments)
		}
	}

	buf.WriteString(fmt.Sprintf("\r\n--%s--", r.Boundary))

	return buf.Bytes(), nil
}

func (r *Request) buildFormData() ([]byte, error) {
	var buf bytes.Buffer

	lenF := len(r.FormData)
	i := 0

	for _, formData := range r.FormData {
		buf.WriteString(fmt.Sprintf("--%s\r\n", r.Boundary))
		buf.WriteString(fmt.Sprintf("%s: %s; name=%q\r\n", HeaderContentDisposition, formData.ContentDisposition, formData.Name))
		buf.WriteString(fmt.Sprintf("%s: %s\r\n\r\n", HeaderContentType, formData.ContentType))
		buf.Write(formData.Body)

		if i < lenF-1 {
			buf.WriteString("\r\n")
		}
		i++
	}

	return buf.Bytes(), nil
}

func (r *Request) buildAttachments() ([]byte, error) {
	var buf bytes.Buffer

	for _, attachment := range r.Attachments {
		buf.WriteString(fmt.Sprintf("--%s\r\n", r.Boundary))
		buf.WriteString(fmt.Sprintf("%s: %s; name=%q; filename=%q\r\n", HeaderContentDisposition, attachment.ContentDisposition, attachment.Name, attachment.FileName))
		buf.WriteString(fmt.Sprintf("%s: %s\r\n\r\n", HeaderContentType, attachment.ContentType))
		buf.Write(attachment.Body)
		buf.WriteString("\r\n")
	}

	return buf.Bytes(), nil
}

func (r *Request) buildZippedAttachments() ([]byte, error) {
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

	for _, attachment := range r.Attachments {
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
