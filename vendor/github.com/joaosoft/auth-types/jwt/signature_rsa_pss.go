package jwt

import (
	"crypto/rand"
	"crypto/rsa"
)

type signatureRSAPSS struct {
	*signatureRSA
	Options *rsa.PSSOptions
}

func (sg *signatureRSAPSS) Verify(headerAndClaims []byte, signature []byte, key interface{}) error {
	var rsaKey *rsa.PublicKey
	switch k := key.(type) {
	case *rsa.PublicKey:
		rsaKey = k
	default:
		return ErrorInvalidAuthorization
	}

	if !sg.Hash.Available() {
		return ErrorInvalidAuthorization
	}
	hasher := sg.Hash.New()
	hasher.Write(headerAndClaims)

	return rsa.VerifyPSS(rsaKey, sg.Hash, hasher.Sum(nil), signature, sg.Options)
}

func (sg *signatureRSAPSS) Signature(headerAndClaims []byte, key interface{}) ([]byte, error) {
	var rsaKey *rsa.PrivateKey

	switch k := key.(type) {
	case *rsa.PrivateKey:
		rsaKey = k
	default:
		return nil, ErrorInvalidAuthorization
	}

	if !sg.Hash.Available() {
		return nil, ErrorInvalidAuthorization
	}

	hasher := sg.Hash.New()
	hasher.Write(headerAndClaims)

	sigBytes, err := rsa.SignPSS(rand.Reader, rsaKey, sg.Hash, hasher.Sum(nil), sg.Options)
	if err != nil {
		return nil, err
	}

	return sigBytes, nil
}
