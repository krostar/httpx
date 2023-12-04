# httpx

[![License](https://img.shields.io/badge/license-MIT-blue)](https://choosealicense.com/licenses/mit/)
![go.mod Go version](https://img.shields.io/github/go-mod/go-version/krostar/httpx?label=go)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/krostar/httpx)
[![Latest tag](https://img.shields.io/github/v/tag/krostar/httpx)](https://github.com/krostar/httpx/tags)
[![Go Report](https://goreportcard.com/badge/github.com/krostar/httpx)](https://goreportcard.com/report/github.com/krostar/httpx)

Useful set of functions and middlewares to use on top of net/http.

## Motivations

On any project where I need a http server, I find myself copy/pasting the same boilerplate over and over again.
This project aims at keeping all that code reusable and properly tested.

```go
func runServer(ctx context.Context) error {
	var (
		address string
		handler http.HandlerFunc
		logger  *slog.Logger
	)

	listener, err := httpx.NewListener(ctx, address)
	if err != nil {
		return fmt.Errorf("unable to create listener: %w", err)
	}

	server := httpx.NewServer(handler, httpx.ServerWithErrorLogger(slog.NewLogLogger(logger.Handler(), slog.LevelWarn)))
	return httpx.Serve(ctx, server, listener, time.Second*15)
}
```

## License

This project is under the MIT licence, please see the LICENCE file.
