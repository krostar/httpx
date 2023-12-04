package httpx

import (
	"log"
)

type serverOptions struct {
	errorLogger *log.Logger
}

// ServerOption defines options applier for the server.
type ServerOption func(*serverOptions)

// ServerWithErrorLogger sets the provider logger to be used for errors in the http server.
func ServerWithErrorLogger(logger *log.Logger) ServerOption {
	return func(o *serverOptions) {
		o.errorLogger = logger
	}
}
