package midx

import (
	"net/http"
)

type monitorOptions struct {
	hasRequestFailed func(r *http.Request) bool
}

// MonitorOption applies option to the http monitor.
type MonitorOption func(*monitorOptions)

// MonitorWithHasRequestFailedFunc sets the function call to check
// whenever the request failed as the provided one.
func MonitorWithHasRequestFailedFunc(fct func(r *http.Request) bool) MonitorOption {
	return func(o *monitorOptions) {
		o.hasRequestFailed = fct
	}
}
