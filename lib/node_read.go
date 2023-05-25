package jd

import (
	"io/ioutil"
	"unicode"

	"github.com/goccy/go-json"
	"gopkg.in/yaml.v2"
)

// ReadJsonFile reads a file as JSON and constructs a JsonNode.
func ReadJsonFile(filename string) (JsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes, json.Unmarshal)
}

// ReadYamlFile reads a file as YAML and constructs a JsonNode.
func ReadYamlFile(filename string) (JsonNode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshal(bytes, yaml.Unmarshal)
}

// ReadJsonString reads a string as JSON and constructs a JsonNode.
func ReadJsonString(s string) (JsonNode, error) {
	return unmarshal([]byte(s), json.Unmarshal)
}

// ReadJsonString reads a string as YAML and constructs a JsonNode.
func ReadYamlString(s string) (JsonNode, error) {
	return unmarshal([]byte(s), yaml.Unmarshal)
}

// ReadJsonBytes reads a byte slice as JSON and constructs a JsonNode.
func ReadJsonBytes(b []byte) (JsonNode, error) {
	return unmarshal(b, json.Unmarshal)
}

// ReadYamlBytes reads a byte slice as YAML and constructs a JsonNode.
func ReadYamlBytes(b []byte) (JsonNode, error) {
	return unmarshal(b, yaml.Unmarshal)
}

func unmarshal(bytes []byte, fn func([]byte, interface{}) error) (JsonNode, error) {
	if isEmptyOrWhitespace(bytes) {
		return voidNode{}, nil
	}
	var v interface{}
	err := fn(bytes, &v)
	if err != nil {
		return nil, err
	}
	n, err := NewJsonNode(v)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func isEmptyOrWhitespace(b []byte) bool {
	for _, v := range b {
		if !unicode.IsSpace(rune(v)) {
			return false
		}
	}
	return true
}
