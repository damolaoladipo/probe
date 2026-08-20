package request

import (
	"fmt"
	"strings"
)

type Request struct {
	Method  string              `yaml:"method"`
	URL     string              `yaml:"url"`
	Headers map[string][]string `yaml:"headers"`
	Body    string              `yaml:"body"`
}

func (r Request) Normalize() Request {
	r.Method = strings.ToUpper(strings.TrimSpace(r.Method))
	if r.Method == "" {
		r.Method = "GET"
	}
	r.URL = strings.TrimSpace(r.URL)
	return r
}

func (r Request) Validate() error {
	if strings.TrimSpace(r.URL) == "" {
		return fmt.Errorf("url is required")
	}
	return nil
}

func ParseHeaders(raw string) (map[string]string, error) {
	out := map[string]string{}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			return nil, fmt.Errorf("header line %q: want Key: Value", line)
		}
		out[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return out, nil
}
