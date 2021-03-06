package common

type Method string

const (
	MethodConnect  Method = "CONNECT"
	MethodGet      Method = "GET"
	MethodHead     Method = "HEAD"
	MethodPost     Method = "POST"
	MethodPut      Method = "PUT"
	MethodPatch    Method = "PATCH"
	MethodDelete   Method = "DELETE"
	MethodOptions  Method = "OPTIONS"
	MethodTrace    Method = "TRACE"
	MethodCopy     Method = "COPY"
	MethodView     Method = "VIEW"
	MethodLink     Method = "LINK"
	MethodUnlink   Method = "UNLINK"
	MethodPurge    Method = "PURGE"
	MethodLock     Method = "LOCK"
	MethodUnlock   Method = "UNLOCK"
	MethodPropFind Method = "PROPFIND"
)

var (
	Methods = []Method{
		MethodGet,
		MethodHead,
		MethodConnect,
		MethodDelete,
		MethodOptions,
		MethodPatch,
		MethodPost,
		MethodTrace,
		MethodPut,
		MethodCopy,
		MethodView,
		MethodLink,
		MethodUnlink,
		MethodPurge,
		MethodLock,
		MethodUnlock,
		MethodPropFind,
	}

	MethodHasBody = map[Method]bool{
		MethodGet:      true,
		MethodHead:     false,
		MethodConnect:  true,
		MethodDelete:   true,
		MethodOptions:  true,
		MethodPatch:    true,
		MethodPost:     true,
		MethodTrace:    true,
		MethodPut:      true,
		MethodCopy:     false,
		MethodView:     true,
		MethodLink:     true,
		MethodUnlink:   true,
		MethodPurge:    false,
		MethodLock:     true,
		MethodUnlock:   false,
		MethodPropFind: true,
	}
)
