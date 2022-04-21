package writers

import (
	"encoding/json"
	"fmt"
)

type FormatHandler func(prefixes map[string]interface{}, tags map[string]interface{}, message interface{}, fields map[string]interface{}, sufixes map[string]interface{}) ([]byte, error)

func JsonFormatHandler(prefixes map[string]interface{}, tags map[string]interface{}, message interface{}, fields map[string]interface{}, sufixes map[string]interface{}) ([]byte, error) {
	if bytes, err := json.Marshal(Message{Prefixes: prefixes, Tags: tags, Message: message, Fields: fields, Sufixes: sufixes}); err != nil {
		return nil, err
	} else {

		return append(bytes, []byte("\n")...), nil
	}
}

func TextFormatHandler(prefixes map[string]interface{}, tags map[string]interface{}, message interface{}, fields map[string]interface{}, sufixes map[string]interface{}) ([]byte, error) {
	type MessageText struct {
		prefixes interface{}
		tags     interface{}
		message  interface{}
		fields   interface{}
		sufixes interface{}
	}

	return []byte(fmt.Sprintf("%+v\n", MessageText{prefixes: prefixes, tags: tags, message: message, fields: fields, sufixes: sufixes})), nil
}
