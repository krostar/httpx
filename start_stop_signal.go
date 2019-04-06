package httpx

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"
)

// StartAndStopWithSignal starts the provided function and calls
// the provided stop function when the provided signals are triggered.
func StartAndStopWithSignal(srv *http.Server, l net.Listener, timeout time.Duration, signals ...os.Signal) error {
	if len(signals) == 0 {
		return errors.Errorf("no signal provided")
	}

	var (
		errChan  = make(chan error, 1)
		stopChan = make(chan os.Signal, 1)
	)
	defer close(stopChan) // make sure we release the goroutine if we can't start
	signal.Notify(stopChan, signals...)

	go func() {
		<-stopChan

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			errChan <- errors.Wrap(err, "unable to gracefully stop")
			return
		}

		errChan <- nil
	}()

	if err := srv.Serve(l); err != nil {
		if err != http.ErrServerClosed {
			return errors.Wrap(err, "unable to start server")
		}
	}

	return <-errChan
}
