package jwt

import (
	"crypto"
	"crypto/hmac"
	"fmt"
)

type signatureHMAC struct {
	Name string
	Hash crypto.Hash
}

func (sg *signatureHMAC) Algorithm() string {
	return sg.Name
}

func (sg *signatureHMAC) Verify(headerAndClaims []byte, signature []byte, key interface{}) error {
	var keyBytes []byte
	switch b := key.(type) {
	case []byte:
		keyBytes = b
	default:
		keyBytes = []byte(fmt.Sprintf("%+v", key))
	}

	if !sg.Hash.Available() {
		return ErrorInvalidAuthorization
	}

	hasher := hmac.New(sg.Hash.New, keyBytes)
	hasher.Write(headerAndClaims)
	if !hmac.Equal(signature, hasher.Sum(nil)) {
		return ErrorInvalidAuthorization
	}

	return nil
}

func (sg *signatureHMAC) Signature(headerAndClaims []byte, key interface{}) ([]byte, error) {
	var keyBytes []byte
	switch b := key.(type) {
	case []byte:
		keyBytes = b
	default:
		keyBytes = []byte(fmt.Sprintf("%+v", key))
	}

	if !sg.Hash.Available() {
		return nil, ErrorInvalidAuthorization
	}

	hasher := hmac.New(sg.Hash.New, keyBytes)
	hasher.Write(headerAndClaims)

	return hasher.Sum(nil), nil
}
