package tui

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "m":
			m.method = nextMethod(m.method)
			return m, nil
		case "enter", "ctrl+enter":
			return m.startSend()
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

	m.urlInput, cmd = m.urlInput.Update(msg)
	return m, cmd
}
