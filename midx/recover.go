package midx

import (
	"net/http"
	"runtime/debug"
)

// Recover is a middleware that recover requests from panic.
func Recover(callbacks ...OnPanicFunc) func(http.Handler) http.Handler {
	callbacks = append(callbacks, defaultPanicCallback)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if reason := recover(); reason != nil {
					for _, callback := range callbacks {
						callback(w, r, reason, debug.Stack())
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
