package tui

import (
	"fmt"
	"strings"
)

// formatHeaders turns response headers into one "Name: value" line each.
// Several values for the same name are joined with ", ". Empty maps render as "(none)".
func formatHeaders(h map[string][]string) string {
	if len(h) == 0 {
		return "(none)\n"
	}
	var b strings.Builder
	for name, values := range h {
		b.WriteString(name)
		b.WriteString(": ")
		b.WriteString(strings.Join(values, ", "))
		b.WriteString("\n")
	}
	return b.String()
}

// View draws the method, URL, editors, status, and response.
// It must not print with fmt.Println; the string returned here is the whole screen.
func (m Model) View() string {
	header := fmt.Sprintf("%-7s %s", m.method, m.urlInput.View())
	s := fmt.Sprintf("Probe\n\n%s\n\nHeaders\n%s\n\nBody\n%s\n\nQuery\n%s\n\n",
		header, m.headers.View(), m.body.View(), m.query.View())
	if m.loading {
		s += "Sending...\n"
	}
	if m.err != nil {
		s += fmt.Sprintf("Error: %v\n", m.err)
	}
	if m.response != nil {
		s += fmt.Sprintf(
			"Response\n\n%s\n%v\n\nHeaders\n%s\nBody\n%s\n",
			m.response.Status,
			m.response.Duration,
			formatHeaders(m.response.Headers),
			m.response.Body,
		)
	}
	s += "\ntab focus   m method (URL)   ctrl+enter send   Ctrl+C quit\n"
	return s
}
