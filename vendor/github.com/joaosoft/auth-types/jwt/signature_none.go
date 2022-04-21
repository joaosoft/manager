package jwt

type signatureNone struct {
	Name string
}

func (sg *signatureNone) Algorithm() string {
	return sg.Name
}

func (sg *signatureNone) Verify(headerAndClaims []byte, signature []byte, key interface{}) (err error) {
	if string(signature) != "" {
		return ErrorInvalidAuthorization
	}

	return nil
}

func (sg *signatureNone) Signature(headerAndClaims []byte, key interface{}) ([]byte, error) {
	return nil, nil
}
