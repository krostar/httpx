package midx

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
)

func TestTracerWithCallback(t *testing.T) {
	var (
		o      Tracer
		called bool
	)

	TracerWithCallback(func(*http.Request, *trace.Span) { called = true })(&o)

	for _, cb := range o.callback {
		cb(nil, nil)
	}

	assert.True(t, called)
}
