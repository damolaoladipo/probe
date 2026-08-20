# Probe

Probe is a testing-first API client for the terminal. Send a request, inspect the
response, then (later) add assertions, save a suite, and run it in CI.

The live app is **v0.0.1**: method + URL + send, with status, headers, body, and
duration on screen.

```text
go run ./cmd/probe
```

Module: `github.com/damola-oladipo/probe`. Go 1.25.7.

## Install

```bash
git clone https://github.com/damola-oladipo/probe
cd probe
go build -o probe ./cmd/probe
./probe
```

## TUI

`probe` opens the workbench (alt screen).

- `m` cycles GET POST PUT PATCH DELETE HEAD OPTIONS
- Enter or `ctrl+enter` sends
- `ctrl+c` quits

Do not use `q` to quit while the URL field is focused; it types into the URL.

## Layout

```text
cmd/probe/              process entry
internal/httpclient/    Send(ctx, request)
internal/request/       method, URL, headers, body
internal/tui/           Model, Update, View
```

The HTTP engine does not import Bubble Tea. The TUI does not import `net/http`.

## Develop

```bash
gofmt -w .
go test ./...
go vet ./...
go run ./cmd/probe
```

See [CONTRIBUTING.md](CONTRIBUTING.md). License: [LICENSE.md](LICENSE.md) (MIT).
