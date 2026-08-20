package tui

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/damola-oladipo/probe/internal/httpclient"
	"github.com/damola-oladipo/probe/internal/request"
)

type Model struct {
	width, height int
	method        string
	urlInput      textinput.Model
	client        *httpclient.Client
	loading       bool
	response      *httpclient.Response
	err           error
}

type responseMsg struct {
	response httpclient.Response
	err      error
}

func NewModel() Model {
	input := textinput.New()
	input.Placeholder = "Enter the API url to test"
	input.Focus()
	input.Width = 60
	return Model{
		urlInput: input,
		client:   httpclient.New(),
		method:   "GET",
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func sendRequest(client *httpclient.Client, req request.Request) tea.Cmd {
	return func() tea.Msg {
		resp, err := client.Send(context.Background(), req)
		return responseMsg{response: resp, err: err}
	}
}

func (m Model) startSend() (Model, tea.Cmd) {
	req := request.Request{
		Method: m.method,
		URL:    strings.TrimSpace(m.urlInput.Value()),
	}.Normalize()
	if err := req.Validate(); err != nil {
		m.err = err
		m.loading = false
		return m, nil
	}
	m.loading = true
	m.err = nil
	m.response = nil
	return m, sendRequest(m.client, req)
}
