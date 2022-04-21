package web

type Encoding string

const (
	EncodingNone     Encoding = "none"
	EncodingChunked  Encoding = "chunked"
	EncodingCompress Encoding = "compress"
	EncodingDeflate  Encoding = "deflate"
	EncodingGzip     Encoding = "gzip"
	EncodingIdentity Encoding = "identity"
)
