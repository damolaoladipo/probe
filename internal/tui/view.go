package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

// pane pads a tab's content to the same width and height as the textareas and viewport.
func (m Model) pane(s string) string {
	w, h := m.paneW, m.paneH
	if w <= 0 {
		w = 60
	}
	if h <= 0 {
		h = 12
	}
	return lipgloss.NewStyle().Width(w).Height(h).MaxHeight(h).Render(s)
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
		body = m.pane(header)
	case tabHeaders:
		body = m.headers.View()
	case tabBody:
		body = m.body.View()
	case tabQuery:
		body = m.query.View()
	case tabResponse:
		if m.err != nil {
			body = m.pane(fmt.Sprintf("Error: %v", m.err))
		} else if m.response != nil {
			body = m.viewport.View()
		} else {
			body = m.pane("No response yet. ctrl+enter to send.")
		}
	}

	s := fmt.Sprintf("Probe\n\n%s\n%s\n%s\n%s", m.renderTabs(), body, m.statusMsg, m.help.View(m.keys))
	if m.width > 0 {
		return m.styles.frame.Width(m.width).Render(s)
	}

	return s
}
