package web

type Protocol string

const (
	ProtocolHttp0p9 Protocol = "HTTP/0.9"
	ProtocolHttp1p0 Protocol = "HTTP/1.0"
	ProtocolHttp1p1 Protocol = "HTTP/1.1"
	ProtocolHttp2   Protocol = "HTTP/2"
)
