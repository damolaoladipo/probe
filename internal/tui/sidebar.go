package tui

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/damola-oladipo/probe/internal/request"
)

const sidebarWidth = 24

// namedRequest is one item in the workbench sidebar.
type namedRequest struct {
	Name string
	Req  request.Request
}

// blankRequest is a new GET with a JSON content type, same as NewModel's editors.
func blankRequest() request.Request {
	return request.Request{
		Method:  "GET",
		Headers: map[string][]string{"Content-Type": {"application/json"}},
	}
}

// labelRequest is the sidebar row: METHOD plus URL path, truncated to fit the pane.
func labelRequest(req request.Request) string {
	req = req.Normalize()
	path := strings.TrimSpace(req.URL)
	if u, err := url.Parse(path); err == nil && u.Path != "" {
		path = u.Path
		if u.RawQuery != "" {
			path += "?" + u.RawQuery
		}
	}
	if path == "" || path == "/" {
		return req.Method + " Untitled"
	}
	if len(path) > 14 {
		path = path[:14] + "..."
	}
	return req.Method + " " + path
}

// stashCurrent copies the editors into the selected sidebar item so switching does not drop edits.
func (m Model) stashCurrent() Model {
	if m.sel < 0 || m.sel >= len(m.requests) {
		return m
	}
	req := m.currentRequest()
	m.requests[m.sel].Req = req
	m.requests[m.sel].Name = labelRequest(req)
	return m
}

// newRequest appends a blank request, selects it, and focuses the sidebar.
func (m Model) newRequest() Model {
	m = m.stashCurrent()
	req := blankRequest()
	m.requests = append(m.requests, namedRequest{Name: "Untitled", Req: req})
	m.sel = len(m.requests) - 1
	m = m.applyLoaded(req)
	m.response = nil
	m.err = nil
	m.statusMsg = "new request"
	m.sidebar = true
	m.urlInput.Blur()
	m.headers.Blur()
	m.body.Blur()
	m.query.Blur()
	return m
}

// selectRequest stashes the current editors, then loads item i into them.
func (m Model) selectRequest(i int) Model {
	if i < 0 || i >= len(m.requests) {
		return m
	}
	m = m.stashCurrent()
	m.sel = i
	m = m.applyLoaded(m.requests[i].Req)
	m.response = nil
	m.err = nil
	return m
}

// deleteRequest removes the selected item. The list always keeps at least one request.
func (m Model) deleteRequest() Model {
	if len(m.requests) <= 1 {
		m.statusMsg = "keep at least one request"
		return m
	}
	m.requests = append(m.requests[:m.sel], m.requests[m.sel+1:]...)
	if m.sel >= len(m.requests) {
		m.sel = len(m.requests) - 1
	}
	m = m.applyLoaded(m.requests[m.sel].Req)
	m.statusMsg = "deleted request"
	return m
}

// renderSidebar draws the named request list to the left of the workbench tabs.
func (m Model) renderSidebar() string {
	h := m.paneH + 3
	if h <= 3 {
		h = 15
	}
	var b strings.Builder
	b.WriteString("Requests\n")
	for i, item := range m.requests {
		name := item.Name
		if name == "" {
			name = labelRequest(item.Req)
		}
		line := fmt.Sprintf("  %s", name)
		if i == m.sel {
			line = "> " + name
		}
		if i == m.sel && m.sidebar {
			line = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Width(sidebarWidth - 2).Render(line)
		} else if i == m.sel {
			line = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229")).Render(line)
		} else {
			line = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(line)
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	border := lipgloss.NormalBorder()
	style := lipgloss.NewStyle().Border(border).BorderForeground(lipgloss.Color("240")).Width(sidebarWidth).Height(h).MaxHeight(h)
	if m.sidebar {
		style = style.BorderForeground(lipgloss.Color("57"))
	}
	return style.Render(b.String())
}
