package writers

type Message struct {
	Prefixes map[string]interface{} `json:"prefixes,omitempty"`
	Tags     map[string]interface{} `json:"tags,omitempty"`
	Message  interface{}            `json:"message,omitempty"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
	Sufixes   map[string]interface{} `json:"sufixes,omitempty"`
}

// IList ...
type IList interface {
	Add(id string, data interface{}) error
	Remove(ids ...string) interface{}
	Size() int
	IsEmpty() bool
	Dump() string
}
