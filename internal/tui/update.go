package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Update is the Bubble Tea event loop.
// Window size resizes the textareas. Keys that belong to the app (quit, send, tab, m)
// are handled here. Everything else is forwarded to the focused editor so Enter
// can be a newline in Headers/Body/Query. Send is ctrl+enter; ctrl+p is the Mac fallback.
// ctrl+s saves request.yaml. ctrl+o loads it. ctrl+n adds a request. ctrl+b focuses the sidebar.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.applySize()
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case responseMsg:
		m.loading = false
		m.cancel = nil
		if msg.err != nil {
			m.err = msg.err
			m.response = nil
			m.viewport.SetContent("")
			return m, nil
		}
		resp := msg.response
		m.response = &resp
		status := m.styles.statusColor(resp.StatusCode).Render(resp.Status)
		m.viewport.SetContent(fmt.Sprintf("%s  %s\n\nHeaders\n%s\nBody\n%s",
			status, resp.Duration, formatHeaders(resp.Headers), prettyJSON(resp.Body)))
		m.viewport.GotoTop()
		return m.setTab(tabResponse)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Send):
			return m.startSend()
		case key.Matches(msg, m.keys.Save):
			m = m.stashCurrent()
			return m.saveRequest()
		case key.Matches(msg, m.keys.Open):
			loaded, cmd := m.loadRequest()
			loaded = loaded.stashCurrent()
			return loaded, cmd
		case key.Matches(msg, m.keys.New):
			return m.newRequest(), nil
		case key.Matches(msg, m.keys.Sidebar):
			m.sidebar = !m.sidebar
			if m.sidebar {
				m.urlInput.Blur()
				m.headers.Blur()
				m.body.Blur()
				m.query.Blur()
				return m, nil
			}
			return m.setTab(m.tab)
		case key.Matches(msg, m.keys.Cancel):
			return m.cancelInFlight()
		case key.Matches(msg, m.keys.NextTab):
			m.sidebar = false
			return m.setTab(nextTab(m.tab))
		case key.Matches(msg, m.keys.PrevTab):
			m.sidebar = false
			return m.setTab(prevTab(m.tab))
		case key.Matches(msg, m.keys.Method) && m.tab == tabRequest && !m.sidebar:
			m.method = nextMethod(m.method)
			return m, nil
		}
		if m.sidebar {
			switch msg.String() {
			case "j", "down":
				return m.selectRequest(m.sel + 1), nil
			case "k", "up":
				return m.selectRequest(m.sel - 1), nil
			case "n":
				return m.newRequest(), nil
			case "d", "x":
				return m.deleteRequest(), nil
			case "enter":
				m.sidebar = false
				return m.setTab(tabRequest)
			}
			return m, nil
		}
	}

	switch m.tab {
	case tabRequest:
		var cmd tea.Cmd
		m.urlInput, cmd = m.urlInput.Update(msg)
		cmds = append(cmds, cmd)
	case tabHeaders:
		var cmd tea.Cmd
		m.headers, cmd = m.headers.Update(msg)
		cmds = append(cmds, cmd)
	case tabBody:
		var cmd tea.Cmd
		m.body, cmd = m.body.Update(msg)
		cmds = append(cmds, cmd)
	case tabQuery:
		var cmd tea.Cmd
		m.query, cmd = m.query.Update(msg)
		cmds = append(cmds, cmd)
	case tabResponse:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
