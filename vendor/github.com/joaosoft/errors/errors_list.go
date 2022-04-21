package errors

import "encoding/json"

func (e *ListErr) Len() int {
	return len(*e)
}

func (e *ListErr) IsEmpty() bool {
	return len(*e) == 0
}

func (e *ListErr) Add(err *Err) *ListErr {
	*e = append(*e, err)
	return e
}

func (e *ListErr) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}
