package request

import (
	"fmt"
	"strings"
)

// Request is the only request type Probe sends.
// Headers use map[string][]string so a name can have more than one value, like net/http.
type Request struct {
	Method  string              `yaml:"method"`
	URL     string              `yaml:"url"`
	Headers map[string][]string `yaml:"headers"`
	Body    string              `yaml:"body"`
}

// Normalize uppercases Method, defaults empty Method to GET, and trims the URL.
// It does not validate; call Validate before Send.
func (r Request) Normalize() Request {

	r.Method = strings.ToUpper(strings.TrimSpace(r.Method))
	if r.Method == "" {
		r.Method = "GET"
	}

	r.URL = strings.TrimSpace(r.URL)
	return r
}

// Validate reports whether the request can be sent.
// The only check today is a non-empty URL after trim.
func (r Request) Validate() error {

	if strings.TrimSpace(r.URL) == "" {
		return fmt.Errorf("url is required")
	}

	return nil
}

// ParseHeaders turns raw textarea lines into request headers.
// Each line must be "Key: Value". Empty lines are skipped.
// Duplicate keys keep every value, in order. A line with no colon is an error.
func ParseHeaders(raw string) (map[string][]string, error) {

	out := map[string][]string{}

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}
		k, v, ok := strings.Cut(line, ":")

		if !ok {
			return nil, fmt.Errorf("header line %q: want Key: Value", line)
		}

		key := strings.TrimSpace(k)
		out[key] = append(out[key], strings.TrimSpace(v))
	}

	return out, nil
}
