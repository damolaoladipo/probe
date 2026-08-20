package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Send    key.Binding
	Quit    key.Binding
	NextTab key.Binding
	PrevTab key.Binding
	Method  key.Binding
	Cancel  key.Binding
}

// newKeyMap is the help-bar bindings.
// Send lists ctrl+enter. Update also accepts ctrl+s because macOS often swallows ctrl+enter.
func newKeyMap() keyMap {
	return keyMap{
		Send:    key.NewBinding(key.WithKeys("ctrl+enter"), key.WithHelp("ctrl+enter", "send")),
		Quit:    key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
		NextTab: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next tab")),
		PrevTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev tab")),
		Method:  key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "method")),
		Cancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

// ShortHelp is the one-line help row at the bottom of the TUI.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Send, k.NextTab, k.Method, k.Cancel, k.Quit}
}

// FullHelp is the expanded help list.
// Until a second row exists, it is one row: the same keys as ShortHelp.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
