package jwt

import (
	"bytes"
	"encoding/base64"
)

type encoderBase64 struct{}

func (e *encoderBase64) Encode(value []byte) ([]byte, error) {
	base64Encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(base64Encoded, value)

	return bytes.TrimRight(base64Encoded, "="), nil
}

func (e *encoderBase64) Decode(value []byte) ([]byte, error) {
	base64Encoded := bytes.NewBuffer(value)
	if l := len(value) % 4; l > 0 {
		base64Encoded.Write(bytes.Repeat([]byte("="), 4-l))
	}

	base64Decoded := make([]byte, base64.URLEncoding.DecodedLen(len(base64Encoded.Bytes())))
	_, err := base64.URLEncoding.Decode(base64Decoded, base64Encoded.Bytes())

	return base64Decoded, err
}
