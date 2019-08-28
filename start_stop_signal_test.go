package httpx

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartAndStopWithSignal(t *testing.T) {
	srv := NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Shutdown(context.Background()) // nolint: errcheck

	listener, err := NewListener(context.Background(), ":0")
	require.NoError(t, err)
	defer listener.Close() // nolint: errcheck

	go func() {
		defer syscall.Kill(os.Getpid(), syscall.SIGUSR1) // nolint: errcheck

		resp, err := http.Get("http://" + listener.Addr().String())
		require.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)
	}()

	require.NoError(t, StartAndStopWithSignal(srv, listener, time.Second, syscall.SIGUSR1))
}

func TestStartAndStopWithSignal_timeout(t *testing.T) {
	srv := NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
	}))
	defer srv.Shutdown(context.Background()) // nolint: errcheck

	listener, err := NewListener(context.Background(), ":0")
	require.NoError(t, err)
	defer listener.Close() // nolint: errcheck

	go func() {
		_, err := http.Get("http://" + listener.Addr().String())
		require.NoError(t, err)
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		syscall.Kill(os.Getpid(), syscall.SIGUSR2) // nolint: errcheck, gosec
	}()

	require.Error(t, StartAndStopWithSignal(srv, listener, time.Millisecond, syscall.SIGUSR2))
}

func TestStartAndStopWithSignal_no_signals(t *testing.T) {
	require.Error(t, StartAndStopWithSignal(nil, nil, 0))
}
