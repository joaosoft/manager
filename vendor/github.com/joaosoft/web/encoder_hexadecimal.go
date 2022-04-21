package web

import (
	"encoding/hex"
)

type encoderHexadecimal struct{}

func (e *encoderHexadecimal) Encode(value []byte) ([]byte, error) {
	hexEncoded := make([]byte, hex.EncodedLen(len(value)))
	n := hex.Encode(hexEncoded, value)

	return hexEncoded[0:n], nil
}

func (e *encoderHexadecimal) Decode(value []byte) ([]byte, error) {
	hexDecoded := make([]byte, hex.DecodedLen(len(value)))
	n, err := hex.Decode(hexDecoded, value)

	return hexDecoded[0:n], err
}
