package jd

import (
	"encoding/json"
	"errors"
	"fmt"
)

type JsonNode interface {
	Json() string
	Equals(n JsonNode) bool
	Diff(n JsonNode) Diff
	diff(n JsonNode, p Path) Diff
	Patch(d Diff) (JsonNode, error)
	patch(pathBehind, pathAhead Path, oldValue, newValue JsonNode) (JsonNode, error)
	// hash(h hash.Hash)
}

func NewJsonNode(n interface{}, options ...option) (JsonNode, error) {
	switch t := n.(type) {
	case map[string]interface{}:
		m := make(jsonObject)
		for k, v := range t {
			if _, ok := v.(JsonNode); !ok {
				e, err := NewJsonNode(v)
				if err != nil {
					return nil, err
				}
				m[k] = e
			}
		}
		return m, nil
	case []interface{}:
		l := make(jsonArray, len(t))
		for i, v := range t {
			if _, ok := v.(JsonNode); !ok {
				e, err := NewJsonNode(v)
				if err != nil {
					return nil, err
				}
				l[i] = e
			}
		}
		return l, nil
	case float64:
		return jsonNumber(t), nil
	case string:
		return jsonString(t), nil
	case bool:
		return jsonBool(t), nil
	case nil:
		return jsonNull{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported type %v", t))
	}
}

func renderJson(n JsonNode) string {
	s, _ := json.Marshal(n)
	// Errors are ignored because JsonNode types are
	// private and known to marshal without error.
	return string(s)
}
