package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	tabRequest = iota
	tabHeaders
	tabBody
	tabQuery
	tabResponse
)

var tabNames = []string{"Request", "Headers", "Body", "Query", "Response"}

// renderTabs draws the active tab highlighted and the rest muted.
// A line of ─ fills the remaining width so the bar looks like a full tab strip.
func (m Model) renderTabs() string {
	active := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Padding(0, 1)
	inactive := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Padding(0, 1)
	parts := make([]string, len(tabNames))
	for i, name := range tabNames {
		if i == m.tab {
			parts[i] = active.Render(name)
		} else {
			parts[i] = inactive.Render(name)
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	line := strings.Repeat("─", max(0, m.paneW-lipgloss.Width(row)))
	return row + line
}

// nextTab wraps forward through Request, Headers, Body, Query, Response.
func nextTab(tab int) int { return (tab + 1) % len(tabNames) }

// prevTab wraps backward through the same list.
func prevTab(tab int) int { return (tab - 1 + len(tabNames)) % len(tabNames) }
