package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
)

type signatureECDSA struct {
	Name      string
	Hash      crypto.Hash
	KeySize   int
	CurveBits int
}

func (sg *signatureECDSA) Algorithm() string {
	return sg.Name
}

func (sg *signatureECDSA) Verify(headerAndClaims []byte, signature []byte, key interface{}) error {
	var ecdsaKey *ecdsa.PublicKey
	switch k := key.(type) {
	case *ecdsa.PublicKey:
		ecdsaKey = k
	default:
		return ErrorInvalidAuthorization
	}

	if len(signature) != 2*sg.KeySize {
		return ErrorInvalidAuthorization
	}

	r := big.NewInt(0).SetBytes(signature[:sg.KeySize])
	s := big.NewInt(0).SetBytes(signature[sg.KeySize:])

	if !sg.Hash.Available() {
	}
	hasher := sg.Hash.New()
	hasher.Write(headerAndClaims)

	if verifystatus := ecdsa.Verify(ecdsaKey, hasher.Sum(nil), r, s); verifystatus == true {
		return nil
	} else {
		return ErrorInvalidAuthorization
	}
}

func (sg *signatureECDSA) Signature(headerAndClaims []byte, key interface{}) ([]byte, error) {
	var ecdsaKey *ecdsa.PrivateKey
	switch k := key.(type) {
	case *ecdsa.PrivateKey:
		ecdsaKey = k
	default:
		return nil, ErrorInvalidAuthorization
	}

	if !sg.Hash.Available() {
		return nil, ErrorInvalidAuthorization
	}

	hasher := sg.Hash.New()
	hasher.Write(headerAndClaims)

	r, s, err := ecdsa.Sign(rand.Reader, ecdsaKey, hasher.Sum(nil))
	if err != nil {
		return nil, err
	}

	curveBits := ecdsaKey.Curve.Params().BitSize

	if sg.CurveBits != curveBits {
		return nil, ErrorInvalidAuthorization
	}

	keyBytes := curveBits / 8
	if curveBits%8 > 0 {
		keyBytes += 1
	}

	rBytes := r.Bytes()
	rBytesPadded := make([]byte, keyBytes)
	copy(rBytesPadded[keyBytes-len(rBytes):], rBytes)

	sBytes := s.Bytes()
	sBytesPadded := make([]byte, keyBytes)
	copy(sBytesPadded[keyBytes-len(sBytes):], sBytes)

	out := append(rBytesPadded, sBytesPadded...)

	return out, nil
}
