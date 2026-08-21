package tui

import (
	"testing"

	"github.com/damola-oladipo/probe/internal/request"
)

func TestLabelRequestUntitled(t *testing.T) {
	got := labelRequest(request.Request{Method: "get"})
	if got != "GET Untitled" {
		t.Fatalf("got %q", got)
	}
}

func TestLabelRequestPath(t *testing.T) {
	got := labelRequest(request.Request{Method: "POST", URL: "https://httpbin.org/post"})
	if got != "POST /post" {
		t.Fatalf("got %q", got)
	}
}

func TestNewRequestKeepsPrevious(t *testing.T) {
	m := NewModel()
	m.urlInput.SetValue("https://example.com/one")
	m = m.newRequest()
	if len(m.requests) != 2 {
		t.Fatalf("len=%d", len(m.requests))
	}
	if m.sel != 1 {
		t.Fatalf("sel=%d", m.sel)
	}
	if m.requests[0].Req.URL != "https://example.com/one" {
		t.Fatalf("stashed url=%q", m.requests[0].Req.URL)
	}
	if m.urlInput.Value() != "" {
		t.Fatalf("new url=%q", m.urlInput.Value())
	}
}

func TestSelectRequestRestores(t *testing.T) {
	m := NewModel()
	m.urlInput.SetValue("https://example.com/one")
	m = m.newRequest()
	m.urlInput.SetValue("https://example.com/two")
	m = m.selectRequest(0)
	if m.urlInput.Value() != "https://example.com/one" {
		t.Fatalf("url=%q", m.urlInput.Value())
	}
	m = m.selectRequest(1)
	if m.urlInput.Value() != "https://example.com/two" {
		t.Fatalf("url=%q", m.urlInput.Value())
	}
}

func TestDeleteKeepsOne(t *testing.T) {
	m := NewModel()
	m = m.deleteRequest()
	if len(m.requests) != 1 {
		t.Fatalf("len=%d", len(m.requests))
	}
}

func TestDeleteSelected(t *testing.T) {
	m := NewModel()
	m.urlInput.SetValue("https://example.com/one")
	m = m.newRequest()
	m = m.deleteRequest()
	if len(m.requests) != 1 {
		t.Fatalf("len=%d", len(m.requests))
	}
	if m.urlInput.Value() != "https://example.com/one" {
		t.Fatalf("url=%q", m.urlInput.Value())
	}
}
