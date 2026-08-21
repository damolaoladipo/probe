package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/damola-oladipo/probe/internal/project"
	"github.com/damola-oladipo/probe/internal/request"
)

// currentRequest builds the request from the editors.
// Parse or query errors are ignored here; startSend is what surfaces them.
func (m Model) currentRequest() request.Request {
	headers, _ := request.ParseHeaders(m.headers.Value())
	u := strings.TrimSpace(m.urlInput.Value())
	u, _ = request.ApplyQuery(u, m.query.Value())
	return request.Request{
		Method:  m.method,
		URL:     u,
		Headers: headers,
		Body:    m.body.Value(),
	}.Normalize()
}

// headersToText turns stored headers back into Key: Value lines for the textarea.
func headersToText(h map[string][]string) string {
	if len(h) == 0 {
		return ""
	}
	var b strings.Builder
	for k, values := range h {
		fmt.Fprintf(&b, "%s: %s\n", k, strings.Join(values, ", "))
	}
	return b.String()
}

// applyLoaded copies a loaded request into the editors.
// Query is cleared; query params live on the URL after ApplyQuery on save.
func (m Model) applyLoaded(req request.Request) Model {
	m.method = req.Method
	m.urlInput.SetValue(req.URL)
	m.body.SetValue(req.Body)
	m.headers.SetValue(headersToText(req.Headers))
	m.query.SetValue("")
	return m
}

// saveRequest overwrites request.yaml in the current working directory.
func (m Model) saveRequest() (Model, tea.Cmd) {
	if err := project.SaveRequest("request.yaml", m.currentRequest()); err != nil {
		m.err = err
	} else {
		m.statusMsg = "saved request.yaml"
		m.err = nil
	}
	return m, nil
}

// loadRequest replaces the editors from request.yaml.
// It does not restore response, loading, or the viewport.
func (m Model) loadRequest() (Model, tea.Cmd) {
	req, err := project.LoadRequest("request.yaml")
	if err != nil {
		m.err = err
		return m, nil
	}
	m = m.applyLoaded(req)
	m.statusMsg = "loaded request.yaml"
	m.err = nil
	return m, nil
}
