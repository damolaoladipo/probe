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
	if h["X-Probe"] != "hello" || h["Accept"] != "json" {
		t.Fatalf("%v", h)
	}
	if _, err := ParseHeaders("no-colon"); err == nil {
		t.Fatal("expected error")
	}
}
