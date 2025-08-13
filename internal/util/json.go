package util

import (
	"bytes"
	"encoding/json"
	"sort"
)

func StableMarshal(v interface{}) ([]byte, error) {
	// First, marshal the struct to JSON to handle the struct type.
	// Then, unmarshal it back into a map[string]interface{} to be processed by the existing stable marshaling logic.
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var asMap map[string]interface{}
	err = json.Unmarshal(bytes, &asMap)
	if err != nil {
		// If it's not a map, just return the marshaled bytes.
		return bytes, nil
	}

	return stableMarshalRec(asMap)
}

func stableMarshalRec(v interface{}) ([]byte, error) {
	switch val := v.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var buf bytes.Buffer
		buf.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			// Marshal key
			keyBytes, err := json.Marshal(k)
			if err != nil {
				return nil, err
			}
			buf.Write(keyBytes)
			buf.WriteByte(':')
			// Marshal value recursively
			valueBytes, err := stableMarshalRec(val[k])
			if err != nil {
				return nil, err
			}
			buf.Write(valueBytes)
		}
		buf.WriteByte('}')
		return buf.Bytes(), nil
	case []interface{}:
		var buf bytes.Buffer
		buf.WriteByte('[')
		for i, item := range val {
			if i > 0 {
				buf.WriteByte(',')
			}
			itemBytes, err := stableMarshalRec(item)
			if err != nil {
				return nil, err
			}
			buf.Write(itemBytes)
		}
		buf.WriteByte(']')
		return buf.Bytes(), nil
	default:
		return json.Marshal(v)
	}
}
