package httpx

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewServer(t *testing.T) {
	srv := NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	}))

	l, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	go func() {
		defer srv.Shutdown(context.Background()) // nolint: errcheck
		resp, err := http.Get("http://" + l.Addr().String())
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
		assert.Equal(t, resp.StatusCode, http.StatusCreated)
	}()

	assert.Equal(t, http.ErrServerClosed, srv.Serve(l))
}

func TestNewServer_with_tls(t *testing.T) {
	srv := NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	}), ServerWithModernTLSConfig())
	srv.Addr = ":7777"

	go func() {
		defer srv.Shutdown(context.Background()) // nolint: errcheck

		rootCAPEM, err := ioutil.ReadFile("./testdata/ca.crt")
		require.NoError(t, err)
		rootCAs := x509.NewCertPool()
		require.True(t, rootCAs.AppendCertsFromPEM(rootCAPEM))

		var client http.Client
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{RootCAs: rootCAs, ServerName: "go-test"}}
		resp, err := client.Get("https://:7777")
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
		require.NotNil(t, resp.TLS)
		assert.Equal(t, resp.StatusCode, http.StatusCreated)
	}()

	assert.Equal(t, http.ErrServerClosed, srv.ListenAndServeTLS("./testdata/cert.crt", "./testdata/cert.key"))
}

func Test_Serve(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	l, err := NewListener(ctx, "localhost:0")
	require.NoError(t, err)
	defer l.Close() // nolint: errcheck

	addr := l.Addr().String()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := Serve(ctx, &http.Server{
			Addr: addr,
			Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusTeapot)
			}),
		}, l, time.Second)
		require.NoError(t, err)
	}()

	resp, err := http.DefaultClient.Get("http://" + addr)
	require.NoError(t, err)
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
	require.NoError(t, resp.Body.Close())

	cancel()

	wg.Wait()
}

func Test_Serve_with_cancellable_context(t *testing.T) {
	ctx, cancel := ContextCancelableBySignal(context.Background(), syscall.SIGUSR1, syscall.SIGUSR2)
	defer cancel()

	l, err := NewListener(ctx, "localhost:0")
	require.NoError(t, err)
	defer l.Close() // nolint: errcheck

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		resp, err := http.Get("http://" + l.Addr().String())
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
	}()

	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
		syscall.Kill(os.Getpid(), syscall.SIGUSR2) // nolint: errcheck, gosec
	}()

	require.NoError(t, Serve(ctx, NewServer(nil), l, time.Millisecond))
	wg.Wait()
}

func Test_ListenAndServe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	err := ListenAndServe(ctx, &http.Server{Addr: "localhost:0"})
	require.NoError(t, err)
}
