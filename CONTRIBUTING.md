# Contributing to Probe

Thank you for working on Probe. This is a Go + Bubble Tea product. Follow the
live code, then this document.

## Before you start

1. Install Go 1.25.7.
2. Do not run `go mod init` again. The module is `github.com/damola-oladipo/probe`.
3. Do not create `internal/app/`. `cmd/probe` wires the TUI.
4. Leave `cmd/e2eprobe` empty until a later version needs a second binary.

## Build and test

```bash
gofmt -w .
go test ./...
go vet ./...
go run ./cmd/probe
```

A change is not done until tests and vet pass. HTTP tests use `httptest`, not the
public internet.

## Architecture

These imports are forbidden:

```text
internal/httpclient  →  internal/tui        NEVER
internal/request     →  internal/tui        NEVER
internal/tui         →  net/http            NEVER
```

- No `fmt.Println` while a Bubble Tea program is running. Errors from `main` go
  to stderr.
- No `go func()` for HTTP. Use `tea.Cmd`.
- `View` only draws. `Send` lives in `internal/httpclient`.
- No new Go interfaces until a second implementation exists. Exception later:
  `assert.Assertion`.
- Stay on Bubble Tea v1 (`github.com/charmbracelet/bubbletea`).
- No emojis in UI, reports, or docs unless a maintainer asks for them.

## How to add a feature

After every change, something must run (`go test` or `go run ./cmd/probe`).

1. Keep packages small. One job per folder. No `pkg/` dump, no `utils/`.
2. Add a package only when a version needs it (`project`, `assert`, `cli`, and so
   on).
3. Functions that wait on the network take `context.Context`. Pure helpers
   (`Normalize`, `View`, `ParseHeaders`) do not.

## Pull requests

- Keep the diff small. One logical change per PR.
- Include tests for `internal/` behavior you change.
- Do not commit `probe` binaries, `debug.log`, or `.env`.
- Do not commit secrets.

## License

By contributing, you agree that your work is licensed under the MIT License in
[LICENSE.md](LICENSE.md).
