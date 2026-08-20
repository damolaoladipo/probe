package tui

import tea "github.com/charmbracelet/bubbletea"

// Update is the Bubble Tea event loop.
// Window size resizes the textareas. Keys that belong to the app (quit, send, tab, m)
// are handled here. Everything else is forwarded to the focused editor so Enter
// can be a newline in Headers/Body/Query. Send is ctrl+enter; ctrl+s is the Mac fallback.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		inner := max(20, msg.Width-4)
		m.headers.SetWidth(inner)
		m.body.SetWidth(inner)
		m.query.SetWidth(inner)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+enter", "ctrl+s":
			return m.startSend()
		case "tab":
			return m.cycleFocus()
		case "m":
			if m.focus == focusURL {
				m.method = nextMethod(m.method)
				return m, nil
			}
		}
	case responseMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			m.response = nil
			return m, nil
		}
		resp := msg.response
		m.response = &resp
		return m, nil
	}

	switch m.focus {
	case focusURL:
		m.urlInput, cmd = m.urlInput.Update(msg)
	case focusHeaders:
		m.headers, cmd = m.headers.Update(msg)
	case focusBody:
		m.body, cmd = m.body.Update(msg)
	case focusQuery:
		m.query, cmd = m.query.Update(msg)
	}
	return m, cmd
}

// cycleFocus moves URL → Headers → Body → Query → URL.
// The old field is blurred and the new one focused so only one caret is active.
func (m Model) cycleFocus() (Model, tea.Cmd) {
	switch m.focus {
	case focusURL:
		m.focus = focusHeaders
		m.urlInput.Blur()
		return m, m.headers.Focus()
	case focusHeaders:
		m.focus = focusBody
		m.headers.Blur()
		return m, m.body.Focus()
	case focusBody:
		m.focus = focusQuery
		m.body.Blur()
		return m, m.query.Focus()
	default:
		m.focus = focusURL
		m.query.Blur()
		return m, m.urlInput.Focus()
	}
}
