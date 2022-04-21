package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

type signatureRSA struct {
	Name string
	Hash crypto.Hash
}

func (sg *signatureRSA) Algorithm() string {
	return sg.Name
}

func (sg *signatureRSA) Verify(headerAndClaims []byte, signature []byte, key interface{}) error {
	var rsaKey *rsa.PublicKey
	var ok bool

	if rsaKey, ok = key.(*rsa.PublicKey); !ok {
		return ErrorInvalidAuthorization
	}

	if !sg.Hash.Available() {
		return ErrorInvalidAuthorization
	}
	hasher := sg.Hash.New()
	hasher.Write(headerAndClaims)

	return rsa.VerifyPKCS1v15(rsaKey, sg.Hash, hasher.Sum(nil), signature)
}

func (sg *signatureRSA) Signature(headerAndClaims []byte, key interface{}) ([]byte, error) {
	var rsaKey *rsa.PrivateKey
	var ok bool

	if rsaKey, ok = key.(*rsa.PrivateKey); !ok {
		return nil, ErrorInvalidAuthorization
	}

	if !sg.Hash.Available() {
		return nil, ErrorInvalidAuthorization
	}

	hasher := sg.Hash.New()
	hasher.Write(headerAndClaims)

	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, sg.Hash, hasher.Sum(nil))
	if err != nil {
		return nil, err
	}

	return sigBytes, nil
}
