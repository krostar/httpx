package httpx

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
)

func Test_NewServer(t *testing.T) {
	srv := NewServer(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	}, ServerWithErrorLogger(nil))

	l, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err)

	go func() {
		resp, err := http.DefaultClient.Get("http://" + l.Addr().String())
		assert.NilError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusCreated)
		assert.NilError(t, srv.Shutdown(context.Background()))
	}()

	assert.ErrorIs(t, srv.Serve(l), http.ErrServerClosed)
}

func Test_Serve(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		l, err := NewListener(ctx, "localhost:0")
		assert.NilError(t, err)
		addr := l.Addr().String()

		var wg errgroup.Group
		wg.Go(func() error {
			srv := NewServer(func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusTeapot)
			})
			return Serve(ctx, srv, l, time.Second)
		})

		resp, err := http.DefaultClient.Get("http://" + addr)
		assert.NilError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusTeapot)

		cancel()
		assert.NilError(t, wg.Wait())
	})

	t.Run("unable to serve", func(t *testing.T) {
		ctx := context.Background()

		srv := NewServer(func(rw http.ResponseWriter, r *http.Request) { rw.WriteHeader(http.StatusTeapot) })

		l, err := NewListener(ctx, "localhost:0")
		assert.NilError(t, err)

		assert.ErrorContains(t, Serve(ctx, srv, listenerFail{Listener: l}, time.Second), "unable to serve listener")
	})

	t.Run("unable to shutdown", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		srv := NewServer(func(rw http.ResponseWriter, r *http.Request) {
			cancel()
			reqCtx, cancelReqCtx := context.WithTimeout(r.Context(), time.Second*2)
			defer cancelReqCtx()
			<-reqCtx.Done()
		})

		l, err := NewListener(context.Background(), "localhost:0")
		assert.NilError(t, err)

		var wg errgroup.Group
		wg.Go(func() error {
			err := Serve(ctx, srv, l, time.Second)
			return err
		})
		wg.Go(func() error {
			_, err := http.DefaultClient.Get("http://" + l.Addr().String())
			return err
		})

		assert.ErrorContains(t, wg.Wait(), "unable to shut server down")
	})
}

func Test_ListenAndServe(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		var wg errgroup.Group
		wg.Go(func() error {
			return ListenAndServe(ctx, "localhost:0", func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusTeapot)
			})
		})

		time.Sleep(time.Second * 2)
		cancel()

		assert.NilError(t, wg.Wait())
	})

	t.Run("ko", func(t *testing.T) {
		ctx := context.Background()

		l, err := NewListener(ctx, "localhost:0")
		assert.NilError(t, err)

		assert.ErrorContains(t, ListenAndServe(ctx, l.Addr().String(), nil), "unable to listen")
	})
}

type listenerFail struct {
	net.Listener
}

func (listenerFail) Accept() (net.Conn, error) { return nil, errors.New("boom") }
