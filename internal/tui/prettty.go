package tui

import (
	"bytes"
	"encoding/json"
)

// prettyJSON indents a JSON body with two spaces.
// If the body is not JSON, the original string is returned unchanged so HTML and text still show.
func prettyJSON(body string) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(body), "", "  "); err != nil {
		return body
	}
	return buf.String()
}
