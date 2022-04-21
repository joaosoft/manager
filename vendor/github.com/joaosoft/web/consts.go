package web

import "time"

const (
	HeaderTimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
	TimeFormat = time.RFC3339
	RegexForURL = "^((http|https)://)?(www)?[a-zA-Z0-9-._:/?&=,%]+$"
)
