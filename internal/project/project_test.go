package project

import (
	"path/filepath"
	"testing"

	"github.com/damola-oladipo/probe/internal/request"
)

// TestSaveRequestRoundTrip writes a request and reads the same method, URL, body, and headers.
func TestSaveRequestRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "request.yaml")
	in := request.Request{
		Method:  "POST",
		URL:     "https://httpbin.org/post",
		Headers: map[string][]string{"Accept": {"application/json"}},
		Body:    `{"n":1}`,
	}
	if err := SaveRequest(path, in); err != nil {
		t.Fatal(err)
	}
	got, err := LoadRequest(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Method != "POST" || got.URL != in.URL || got.Body != in.Body {
		t.Fatalf("%+v", got)
	}
	if len(got.Headers["Accept"]) != 1 || got.Headers["Accept"][0] != "application/json" {
		t.Fatalf("headers %+v", got.Headers)
	}
}

// TestLoadRequestMissingFile errors when the path does not exist.
func TestLoadRequestMissingFile(t *testing.T) {
	_, err := LoadRequest("does-not-exist.yaml")
	if err == nil {
		t.Fatal("expected error")
	}
}

// TestLoadRequestFixture loads testdata/get.yaml, including a scalar header value.
func TestLoadRequestFixture(t *testing.T) {
	got, err := LoadRequest("testdata/get.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if got.Method != "GET" || got.URL == "" {
		t.Fatalf("%+v", got)
	}
}
