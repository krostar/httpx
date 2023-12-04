package httpx

import "net/http"

// MiddlewareRecoverFunc is the signature of the callback called each time the request panics.
type MiddlewareRecoverFunc func(rw http.ResponseWriter, r *http.Request, reason any)

// MiddlewareRecover is a middleware that recover requests from a panic.
func MiddlewareRecover(callbacks ...MiddlewareRecoverFunc) func(http.Handler) http.Handler {
	if len(callbacks) == 0 {
		callbacks = append(callbacks, func(rw http.ResponseWriter, _ *http.Request, _ any) {
			rw.WriteHeader(http.StatusInternalServerError)
		})
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if reason := recover(); reason != nil {
					for _, callback := range callbacks {
						callback(w, r, reason)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
