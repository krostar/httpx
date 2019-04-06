package midx

import "net/http"

// OnPanicFunc is the signature of the callback called each time the request panics.
type OnPanicFunc func(w http.ResponseWriter, r *http.Request, reason interface{}, stack []byte)

func defaultPanicCallback(w http.ResponseWriter, _ *http.Request, _ interface{}, _ []byte) {
	w.WriteHeader(http.StatusInternalServerError)
}
