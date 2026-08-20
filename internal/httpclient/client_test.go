package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/damola-oladipo/probe/internal/request"
)

func TestClientSendGET(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))

	defer srv.Close()

	resp, err := New().Send(context.Background(), request.Request{Method: http.MethodGet, URL: srv.URL})
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 || resp.Body != `{"ok":true}` {
		t.Fatalf("%d %q", resp.StatusCode, resp.Body)
	}
}

func TestClientSendPOSTBody(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, _ := io.ReadAll(r.Body)
		if string(got) != `{"n":1}` {
			t.Fatalf("body = %q", got)
		}
		w.WriteHeader(http.StatusCreated)
	}))

	defer srv.Close()

	resp, err := New().Send(context.Background(), request.Request{
		Method: http.MethodPost, URL: srv.URL, Body: `{"n":1}`,
	})

	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("%v %d", err, resp.StatusCode)
	}
}

func TestClientSendHeaders(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Probe") != "yes" {
			t.Fatalf("X-Probe = %q", r.Header.Get("X-Probe"))
		}
		w.WriteHeader(http.StatusOK)
	}))

	defer srv.Close()

	_, err := New().Send(context.Background(), request.Request{
		Method: http.MethodGet, URL: srv.URL, Headers: map[string][]string{"X-Probe": {"yes"}},
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestClientSendHEADNoBody(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, _ := io.ReadAll(r.Body)
		if len(got) != 0 {
			t.Fatalf("HEAD sent a body: %q", got)
		}
		w.WriteHeader(http.StatusOK)
	}))

	defer srv.Close()
	_, err := New().Send(context.Background(), request.Request{
		Method: http.MethodHead, URL: srv.URL, Body: "should-not-be-sent",
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestClientSendInvalidURL(t *testing.T) {

	_, err := New().Send(context.Background(), request.Request{Method: http.MethodGet, URL: "://not-a-url"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestClientSendBodyTooLarge(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.Copy(w, strings.NewReader(strings.Repeat("a", maxBodyBytes+1)))
	}))

	defer srv.Close()

	resp, err := New().Send(context.Background(), request.Request{Method: http.MethodGet, URL: srv.URL})
	if err == nil || err.Error() != "response body exceeds 8MiB" {
		t.Fatalf("err = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status should still be set, got %d", resp.StatusCode)
	}
}
