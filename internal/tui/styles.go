package tui

import "github.com/charmbracelet/lipgloss"

type styles struct {
	frame  lipgloss.Style
	ok     lipgloss.Style
	redir  lipgloss.Style
	client lipgloss.Style
	server lipgloss.Style
}

// newStyles is the Lip Gloss chrome: rounded frame plus status colors.
// Colors follow HTTP classes: 2xx green, 3xx cyan, 4xx yellow, 5xx red.
func newStyles() styles {
	return styles{
		frame:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1),
		ok:     lipgloss.NewStyle().Foreground(lipgloss.Color("42")),
		redir:  lipgloss.NewStyle().Foreground(lipgloss.Color("51")),
		client: lipgloss.NewStyle().Foreground(lipgloss.Color("220")),
		server: lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
	}
}

// statusColor picks a style from the HTTP status class.
// Codes outside 2xx–5xx get an unstyled lipgloss.Style.
func (s styles) statusColor(code int) lipgloss.Style {
	switch {
	case code >= 200 && code < 300:
		return s.ok
	case code >= 300 && code < 400:
		return s.redir
	case code >= 400 && code < 500:
		return s.client
	case code >= 500:
		return s.server
	default:
		return lipgloss.NewStyle()
	}
}
