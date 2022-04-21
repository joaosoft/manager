package web

type Schema string

const (
	SchemaNone  Schema = ""
	SchemaHttp  Schema = "http"
	SchemaHttps Schema = "https"
)
