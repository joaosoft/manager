package errors

type ListErr []*Err

type Err struct {
	Previous *Err   `json:"previous,omitempty"`
	Level    Level  `json:"level"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Stack    string `json:"stack"`
}
