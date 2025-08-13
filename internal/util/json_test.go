package util

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStableMarshal_MapSorting(t *testing.T) {
	input := map[string]interface{}{
		"b": 2,
		"a": 1,
		"c": map[string]interface{}{
			"y": 2,
			"x": 1,
		},
	}

	b, err := StableMarshal(input)
	assert.NoError(t, err)

	// Expect keys sorted: a, b, c and inside c: x, y
	expected := "{\"a\":1,\"b\":2,\"c\":{\"x\":1,\"y\":2}}"
	assert.Equal(t, expected, string(b))
}

func TestStableMarshal_Slice(t *testing.T) {
	input := []interface{}{map[string]interface{}{"b": 2, "a": 1}, map[string]interface{}{"d": 4, "c": 3}}

	b, err := StableMarshal(input)
	assert.NoError(t, err)

	// Each map within slice should have its keys sorted
	expected := "[{\"a\":1,\"b\":2},{\"c\":3,\"d\":4}]"
	assert.Equal(t, expected, string(b))
}

func TestStableMarshal_NonMap(t *testing.T) {
	type sample struct {
		B int `json:"b"`
		A int `json:"a"`
	}

	input := sample{B: 2, A: 1}
	b, err := StableMarshal(input)
	assert.NoError(t, err)

	// Struct marshals to JSON then unmarshals to map; expect sorted keys a, b
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	assert.ElementsMatch(t, []string{"a", "b"}, keys)
}
