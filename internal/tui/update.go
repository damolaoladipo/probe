package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Update is the Bubble Tea event loop.
// Window size resizes the textareas. Keys that belong to the app (quit, send, tab, m)
// are handled here. Everything else is forwarded to the focused editor so Enter
// can be a newline in Headers/Body/Query. Send is ctrl+enter; ctrl+s is the Mac fallback.
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
		m.viewport.SetContent(prettyJSON(resp.Body))
		m.viewport.GotoTop()
		return m.setTab(tabResponse)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Send):
			return m.startSend()
		case key.Matches(msg, m.keys.Cancel):
			return m.cancelInFlight()
		case key.Matches(msg, m.keys.NextTab):
			return m.setTab(nextTab(m.tab))
		case key.Matches(msg, m.keys.PrevTab):
			return m.setTab(prevTab(m.tab))
		case key.Matches(msg, m.keys.Method) && m.tab == tabRequest:
			m.method = nextMethod(m.method)
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

