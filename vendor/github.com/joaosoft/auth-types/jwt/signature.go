package jwt

import (
	"crypto"
	"crypto/rsa"
)

const (
	SignatureES256 signature = "ES256"
	SignatureES384 signature = "ES384"
	SignatureES512 signature = "ES512"
	SignatureRS256 signature = "RS256"
	SignatureRS384 signature = "RS384"
	SignatureRS512 signature = "RS512"
	SignaturePS256 signature = "PS256"
	SignaturePS384 signature = "PS384"
	SignaturePS512 signature = "PS512"
	SignatureHS256 signature = "HS256"
	SignatureHS384 signature = "HS384"
	SignatureHS512 signature = "HS512"
	SignatureNONE  signature = "NONE"
)

type signature string
type isignature interface {
	Verify(headerAndClaims []byte, signature []byte, key interface{}) error
	Signature(headerAndClaims []byte, key interface{}) ([]byte, error)
	Algorithm() string
}

var signatureMethods = map[signature]isignature{
	SignatureES256: &signatureECDSA{Name: string(SignatureES256), Hash: crypto.SHA256, KeySize: 32, CurveBits: 256},
	SignatureES384: &signatureECDSA{Name: string(SignatureES384), Hash: crypto.SHA384, KeySize: 48, CurveBits: 384},
	SignatureES512: &signatureECDSA{Name: string(SignatureES512), Hash: crypto.SHA512, KeySize: 66, CurveBits: 521},
	SignatureRS256: &signatureRSA{Name: string(SignatureRS256), Hash: crypto.SHA256},
	SignatureRS384: &signatureRSA{Name: string(SignatureRS384), Hash: crypto.SHA384},
	SignatureRS512: &signatureRSA{Name: string(SignatureRS512), Hash: crypto.SHA512},
	SignaturePS256: &signatureRSAPSS{signatureRSA: &signatureRSA{Name: string(SignaturePS256), Hash: crypto.SHA256}, Options: &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA256}},
	SignaturePS384: &signatureRSAPSS{signatureRSA: &signatureRSA{Name: string(SignaturePS384), Hash: crypto.SHA384}, Options: &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA384}},
	SignaturePS512: &signatureRSAPSS{signatureRSA: &signatureRSA{Name: string(SignaturePS512), Hash: crypto.SHA512}, Options: &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA512}},
	SignatureHS256: &signatureHMAC{string(SignatureHS256), crypto.SHA256},
	SignatureHS384: &signatureHMAC{string(SignatureHS384), crypto.SHA384},
	SignatureHS512: &signatureHMAC{string(SignatureHS512), crypto.SHA512},
	SignatureNONE:  &signatureNone{Name: string(SignatureNONE)},
}
