package request

import "testing"

func TestValidateEmptyURL(t *testing.T) {
	if (Request{Method: "GET"}).Validate() == nil {
		t.Fatal("expected error")
	}
}

func TestParseHeaders(t *testing.T) {
	h, err := ParseHeaders("X-Probe: hello\n\nAccept: json")
	if err != nil {
		t.Fatal(err)
	}
	if len(h["X-Probe"]) != 1 || h["X-Probe"][0] != "hello" ||
		len(h["Accept"]) != 1 || h["Accept"][0] != "json" {
		t.Fatalf("%v", h)
	}
	if _, err := ParseHeaders("no-colon"); err == nil {
		t.Fatal("expected error")
	}
}

func TestApplyQuery(t *testing.T) {
	got, err := ApplyQuery("http://h/path", "a=1\nq=probe")
	if err != nil {
		t.Fatal(err)
	}
	if got != "http://h/path?a=1&q=probe" && got != "http://h/path?q=probe&a=1" {
		t.Fatalf("got %q", got)
	}
	same, err := ApplyQuery("http://h/path?x=1", "")
	if err != nil || same != "http://h/path?x=1" {
		t.Fatalf("%q %v", same, err)
	}
	if _, err := ApplyQuery("http://h/path", "nocolon"); err == nil {
		t.Fatal("expected error")
	}
}
