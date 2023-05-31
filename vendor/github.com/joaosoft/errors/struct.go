package errors

type ErrorList []*Error

type Error struct {
	Previous *Error      `json:"previous,omitempty"`
	Level    Level       `json:"level"`
	Code     interface{} `json:"code"`
	Message  string      `json:"message"`
	Stack    string      `json:"stack"`
}
