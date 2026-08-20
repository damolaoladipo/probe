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
	if m.loading {
		header += "  " + m.spinner.View()
	}
	
	body := ""
	switch m.tab {
	case tabRequest:
		body = header
	case tabHeaders:
		body = m.headers.View()
	case tabBody:
		body = m.body.View()
	case tabQuery:
		body = m.query.View()
	case tabResponse:
		if m.err != nil {
			body = fmt.Sprintf("Error: %v", m.err)
		} else if m.response != nil {
			status := m.styles.statusColor(m.response.StatusCode).Render(m.response.Status)
			body = fmt.Sprintf("%s  %s\n\nHeaders\n%s\nBody\n%s",
				status, m.response.Duration, formatHeaders(m.response.Headers), m.viewport.View())
		} else {
			body = "No response yet. ctrl+enter to send."
		}
	}
	
	s := fmt.Sprintf("Probe\n%s\n%s\n%s", m.renderTabs(), body, m.help.View(m.keys))
	if m.width > 0 {
		return m.styles.frame.Width(m.width).Render(s)
	}
	
	return s
}