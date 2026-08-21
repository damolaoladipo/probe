package project

import (
	"fmt"
	"os"

	"github.com/damola-oladipo/probe/internal/request"
	"gopkg.in/yaml.v3"
)

// SaveRequest writes req as YAML, overwriting path.
// Permission is 0o644. This package must not import bubbletea.
func SaveRequest(path string, req request.Request) error {
	raw, err := yaml.Marshal(req)
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}

// LoadRequest reads a YAML request file and returns a normalized Request.
// Headers may be a string (Accept: application/json) or a list; both become []string.
func LoadRequest(path string) (request.Request, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return request.Request{}, err
	}

	// Decode headers as yaml.Node so a scalar and a sequence both work.
	var file struct {
		Method  string               `yaml:"method"`
		URL     string               `yaml:"url"`
		Headers map[string]yaml.Node `yaml:"headers"`
		Body    string               `yaml:"body"`
	}
	if err := yaml.Unmarshal(raw, &file); err != nil {
		return request.Request{}, err
	}

	headers := map[string][]string{}
	for k, node := range file.Headers {
		vals, err := headerValues(node)
		if err != nil {
			return request.Request{}, fmt.Errorf("header %s: %w", k, err)
		}
		headers[k] = vals
	}

	return request.Request{
		Method:  file.Method,
		URL:     file.URL,
		Headers: headers,
		Body:    file.Body,
	}.Normalize(), nil
}

// headerValues turns one YAML header value into []string.
// A scalar is one value. A sequence is many. Kind 0 is an empty node.
func headerValues(n yaml.Node) ([]string, error) {
	switch n.Kind {
	case yaml.ScalarNode:
		return []string{n.Value}, nil
	case yaml.SequenceNode:
		var out []string
		if err := n.Decode(&out); err != nil {
			return nil, err
		}
		return out, nil
	case 0:
		return nil, nil
	default:
		return nil, fmt.Errorf("want string or list, got kind %d", n.Kind)
	}
}
