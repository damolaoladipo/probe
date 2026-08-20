package tui

import "testing"

// TestPrettyJSONFallback leaves non-JSON bodies unchanged.
func TestPrettyJSONFallback(t *testing.T) {
	raw := "<html>nope</html>"
	if prettyJSON(raw) != raw {
		t.Fatal("expected raw fallback")
	}
}
