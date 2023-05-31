package jwt

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/joaosoft/errors"
)

var (
	ErrorInvalidAuthorization   = errors.New(errors.LevelError, http.StatusUnauthorized, "invalid authorization")
	ErrorInvalidSignatureMethod = errors.New(errors.LevelError, http.StatusUnauthorized, "invalid signature method")
	ErrorInvalidJwtAlgorithm    = errors.New(errors.LevelError, http.StatusUnauthorized, "invalid signature method")
	ErrorClaimsValidation       = errors.New(errors.LevelError, http.StatusUnauthorized, "error on claims validation")
)

type KeyFunc func(*Token) (interface{}, error)
type CheckFunc func(Claims) (bool, error)

type Token struct {
	raw string `json:"raw"`

	headers   map[string]interface{} `json:"headers"`
	payload   string                 `json:"payload"`
	signature string                 `json:"signature"`

	claims   Claims
	method   isignature
	encoders []iencoder
}

func New(signature signature) *Token {
	method := signatureMethods[signature]

	return &Token{
		headers: map[string]interface{}{
			HeaderTypeKey:      constHeaderTypeJwt,
			HeaderAlgorithmKey: method.Algorithm(),
		},
		claims:   Claims{},
		method:   method,
		encoders: []iencoder{encoderMethods[EncodeBase64]},
	}
}

func (t *Token) Generate(claims Claims, key interface{}) (string, error) {
	t.claims = claims

	headersMarshal, err := json.Marshal(t.headers)
	if err != nil {
		return "", err
	}

	claimsMarshal, err := json.Marshal(t.claims)
	if err != nil {
		return "", err
	}

	headersEncoded, err := t.encode(headersMarshal)
	if err != nil {
		return "", err
	}

	claimsEncoded, err := t.encode(claimsMarshal)
	if err != nil {
		return "", err
	}

	headerAndClaims := strings.Join([]string{headersEncoded, claimsEncoded}, ".")

	signature, err := t.method.Signature([]byte(headerAndClaims), key)
	if err != nil {
		return "", err
	}

	signatureEncoded, err := t.encode(signature)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{string(headerAndClaims), signatureEncoded}, "."), nil
}

func Check(tokenString string, keyFunc KeyFunc, checkFunc CheckFunc, claims Claims, skipClaims bool) (bool, error) {
	token := &Token{
		raw:      tokenString,
		headers:  make(map[string]interface{}),
		claims:   claims,
		encoders: []iencoder{encoderMethods[EncodeBase64]},
	}

	split := strings.Split(tokenString, ".")
	if len(split) != 3 {
		return false, nil
	}

	// headers
	decodedHeader, err := token.decode(split[0])
	if err != nil {
		return false, err
	}

	if err = json.Unmarshal(decodedHeader, &token.headers); err != nil {
		return false, err
	}

	// Claims
	token.claims = claims

	decodedClaims, err := token.decode(split[1])
	if err != nil {
		return false, err

	}

	if err = json.Unmarshal(decodedClaims, &token.claims); err != nil {
		return false, err
	}

	// signature
	if method, ok := token.headers[HeaderAlgorithmKey].(string); ok {
		if token.method, ok = signatureMethods[signature(method)]; !ok {
			return false, ErrorInvalidSignatureMethod
		}
	} else {
		return false, ErrorInvalidJwtAlgorithm
	}

	// execute keyFunc to get the key
	key, err := keyFunc(token)
	if err != nil {
		return false, err
	}

	// Claims
	if !skipClaims {
		if !claims.Validate() {
			return false, ErrorClaimsValidation
		}
	}

	// signature validation
	sig, err := token.decode(split[2])
	if err != nil {
		return false, err
	}

	token.signature = string(sig)

	// header.claims
	decodedHeaderAndClaims := strings.Join(split[0:2], ".")

	if err = token.method.Verify([]byte(decodedHeaderAndClaims), sig, key); err != nil {
		return false, err
	}

	// check claims
	return checkFunc(token.claims)
}

func (t *Token) encode(seg []byte) (string, error) {
	var err error
	next := seg

	for _, encoder := range t.encoders {
		next, err = encoder.Encode(next)
		if err != nil {
			return "", err
		}
	}

	return string(next), nil
}

func (t *Token) decode(seg string) ([]byte, error) {
	var err error
	next := []byte(seg)
	lenE := len(t.encoders) - 1

	for i, _ := range t.encoders {
		next, err = t.encoders[lenE-i].Decode(next)
		if err != nil {
			return nil, err
		}
	}

	return next, nil
}
