package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

func Test_MiddlewareRecover(t *testing.T) {
	handler := func(http.ResponseWriter, *http.Request) { panic("eww") }

	r := httptest.NewRequest("", "/", nil)

	t.Run("defaults", func(t *testing.T) {
		rw := httptest.NewRecorder()

		MiddlewareRecover()(http.HandlerFunc(handler)).ServeHTTP(rw, r)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})

	t.Run("with callbacks", func(t *testing.T) {
		rw := httptest.NewRecorder()

		MiddlewareRecover(func(rw http.ResponseWriter, r *http.Request, reason any) {
			assert.Equal(t, "eww", reason)
			rw.WriteHeader(http.StatusServiceUnavailable)
		})(http.HandlerFunc(handler)).ServeHTTP(rw, r)

		assert.Equal(t, http.StatusServiceUnavailable, rw.Code)
	})
}
