package httpx

import (
	"context"
	"os"
	"os/signal"
)

// ContextCancelableBySignal cancels a context by provided a signal.
func ContextCancelableBySignal(ctx context.Context, signals ...os.Signal) (context.Context, func()) {
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	clean := func() {
		signal.Ignore(signals...)
		close(signalChan)
	}

	// catch some stop signals, and cancel the context if caught
	signal.Notify(signalChan, signals...)
	go func() {
		<-signalChan // block until a signal is received
		cancel()
	}()

	return ctx, clean
}
