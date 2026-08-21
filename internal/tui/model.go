package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/damola-oladipo/probe/internal/httpclient"
	"github.com/damola-oladipo/probe/internal/request"
)

const (
	focusURL = iota
	focusHeaders
	focusBody
	focusQuery
)

// Model is the Bubble Tea workbench: editors, inspect pane, and in-flight HTTP.
// statusMsg is the last save/load line drawn above the help bar.
type Model struct {
	width, height int
	paneW, paneH  int
	tab           int
	method        string
	focus         int
	urlInput      textinput.Model
	headers       textarea.Model
	body          textarea.Model
	query         textarea.Model
	viewport      viewport.Model
	spinner       spinner.Model
	help          help.Model
	keys          keyMap
	styles        styles
	client        *httpclient.Client
	loading       bool
	response      *httpclient.Response
	err           error
	statusMsg     string
	cancel        context.CancelFunc
}

type responseMsg struct {
	response httpclient.Response
	err      error
}

// NewModel builds the workbench: URL, headers, body, query, spinner, viewport, and HTTP client.
// The URL field starts focused. Method is GET. Editors other than URL start blurred
// so tab can move focus without two carets.
func NewModel() Model {

	// Create the headers textarea
	headers := textarea.New()
	headers.Placeholder = "Content-Type: application/json"
	headers.SetWidth(60)
	headers.SetHeight(12)
	headers.SetValue("Content-Type: application/json")
	headers.Blur()

	// Create the URL input
	input := textinput.New()
	input.Placeholder = "Enter the API url to test"
	input.Focus()
	input.Width = 60

	// Create the body textarea
	body := textarea.New()
	body.Placeholder = `{"name":"Damola"}`
	body.SetWidth(60)
	body.SetHeight(12)
	body.Blur()

	// Create the query textarea
	query := textarea.New()
	query.Placeholder = "page=2"
	query.SetWidth(60)
	query.SetHeight(12)
	query.Blur()

	// Create the sending spinner
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	// Create the response viewport
	vp := viewport.New(80, 12)

	m := Model{
		urlInput: input,
		headers:  headers,
		body:     body,
		query:    query,
		viewport: vp,
		spinner:  sp,
		help:     help.New(),
		keys:     newKeyMap(),
		styles:   newStyles(),
		client:   httpclient.New(),
		method:   "GET",
		tab:      tabRequest,
	}

	return m
}

// Init starts commands that should run as soon as the program launches.
// Blink is the URL cursor. Tick starts the spinner so it is ready when a send begins.
func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

// startSend parses the form, cancels any in-flight request, and sends.
// Bad header lines (no ":") or query lines (no "=") set m.err and do not send.
// The request runs in a tea.Cmd with a cancelable context so esc/a new send can abort it.
func (m Model) startSend() (Model, tea.Cmd) {

	// Parse Key: Value lines from the headers textarea
	headers, err := request.ParseHeaders(m.headers.Value())
	if err != nil {
		m.err = err
		m.loading = false
		return m, nil
	}

	// Apply query lines onto the URL
	u, err := request.ApplyQuery(strings.TrimSpace(m.urlInput.Value()), m.query.Value())
	if err != nil {
		m.err = err
		m.loading = false
		return m, nil
	}

	// Build the request from the current editors
	req := request.Request{
		Method:  m.method,
		URL:     u,
		Headers: headers,
		Body:    m.body.Value(),
	}.Normalize()

	if err := req.Validate(); err != nil {
		m.err = err
		m.loading = false
		return m, nil
	}

	// Cancel the in-flight request if one is running.
	// The previous tea.Cmd still finishes, but Send should see ctx.Done().
	if m.cancel != nil {
		m.cancel()
	}

	// Create a new context with cancellation.
	// Store cancel on the model so the next send or esc can stop this one.
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	m.loading = true
	m.err = nil
	m.statusMsg = ""
	m.response = nil

	// Send on a tea.Cmd so the TUI stays responsive.
	// Bubble Tea runs this function; the TUI must not use go func() for HTTP.
	return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
		resp, err := m.client.Send(ctx, req)
		return responseMsg{response: resp, err: err}
	})

}

func (m Model) cancelInFlight() (Model, tea.Cmd) {

	if m.loading && m.cancel != nil {
		m.cancel()
		m.cancel = nil
		m.loading = false
		m.err = fmt.Errorf("request cancelled")
	}
	return m, nil
}

func (m *Model) applySize() {
	// One pane size for every tab so switching Request/Headers/Body/Query/Response does not jump.
	inner := max(20, m.width-4)
	paneH := max(8, m.height-10)
	m.paneW = inner
	m.paneH = paneH
	m.urlInput.Width = max(10, inner-10)
	m.headers.SetWidth(inner)
	m.body.SetWidth(inner)
	m.query.SetWidth(inner)
	m.headers.SetHeight(paneH)
	m.body.SetHeight(paneH)
	m.query.SetHeight(paneH)
	m.viewport.Width = inner
	m.viewport.Height = paneH
}

func (m Model) setTab(tab int) (Model, tea.Cmd) {
	m.tab = tab
	m.urlInput.Blur()
	m.headers.Blur()
	m.body.Blur()
	m.query.Blur()
	switch tab {
	case tabRequest:
		return m, m.urlInput.Focus()
	case tabHeaders:
		return m, m.headers.Focus()
	case tabBody:
		return m, m.body.Focus()
	case tabQuery:
		return m, m.query.Focus()
	default:
		return m, nil
	}
}
