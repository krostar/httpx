package midx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPanic_withDefaultCallback(t *testing.T) {
	var (
		handler = func(w http.ResponseWriter, r *http.Request) {
			panic("eww.")
		}
		r = httptest.NewRequest("", "/", nil)
		w = httptest.NewRecorder()
	)

	require.NotPanics(t, func() {
		Recover()(http.HandlerFunc(handler)).ServeHTTP(w, r)
	}, "handler paniced but middleware should have handled it")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPanic_withCallback(t *testing.T) {
	var (
		handler = func(w http.ResponseWriter, r *http.Request) {
			panic("eww.")
		}
		r = httptest.NewRequest("", "/", nil)
		w = httptest.NewRecorder()
	)

	require.NotPanics(t, func() {
		Recover(func(w http.ResponseWriter, _ *http.Request, reason interface{}, _ []byte) {
			assert.NotNil(t, reason, "if reason is nil it's means there was not panic")
			w.WriteHeader(http.StatusServiceUnavailable)
		})(http.HandlerFunc(handler)).ServeHTTP(w, r)
	}, "handler paniced but middleware should have handled it")

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}
