package tui

import "github.com/charmbracelet/bubbles/key"

// keyMap is the help-bar and Update bindings for the workbench.
type keyMap struct {
	Send    key.Binding
	Save    key.Binding
	Open    key.Binding
	New     key.Binding
	Sidebar key.Binding
	Quit    key.Binding
	NextTab key.Binding
	PrevTab key.Binding
	Method  key.Binding
	Cancel  key.Binding
}

// newKeyMap is the help-bar bindings.
// Send is ctrl+enter. ctrl+p is the Mac fallback because terminals swallow ctrl+enter.
// ctrl+s is save (learn-004), not send.
func newKeyMap() keyMap {
	return keyMap{
		Send:    key.NewBinding(key.WithKeys("ctrl+enter", "ctrl+p"), key.WithHelp("ctrl+enter", "send")),
		Save:    key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
		Open:    key.NewBinding(key.WithKeys("ctrl+o"), key.WithHelp("ctrl+o", "open")),
		New:     key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "new")),
		Sidebar: key.NewBinding(key.WithKeys("ctrl+b"), key.WithHelp("ctrl+b", "sidebar")),
		Quit:    key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
		NextTab: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next tab")),
		PrevTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev tab")),
		Method:  key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "method")),
		Cancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

// ShortHelp is the one-line help row at the bottom of the TUI.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Send, k.New, k.Save, k.Open, k.Sidebar, k.NextTab, k.Method, k.Cancel, k.Quit}
}

// FullHelp is the expanded help list.
// Until a second row exists, it is one row: the same keys as ShortHelp.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
