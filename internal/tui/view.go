package tui

import (
	"fmt"
	"strings"
)

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

func (m Model) View() string {
	header := fmt.Sprintf("%-7s %s", m.method, m.urlInput.View())
	s := fmt.Sprintf("Probe\n\n%s\n\n", header)
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
	s += "\nm method   Enter/ctrl+enter send   Ctrl+C quit\n"
	return s
}
