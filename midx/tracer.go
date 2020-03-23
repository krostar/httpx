package midx

import (
	"net/http"

	"github.com/krostar/httpinfo"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
)

// Tracer stores the configuration for the tracer.
type Tracer struct {
	propagation propagation.HTTPFormat
	sampler     trace.Sampler
	callback    []func(r *http.Request, span *trace.Span)
}

// NewTracer creates a new tracer with the provided options.
func NewTracer(opts ...TracerOption) *Tracer {
	tracer := Tracer{
		propagation: &b3.HTTPFormat{},
		sampler:     trace.AlwaysSample(),
	}

	for _, opt := range opts {
		opt(&tracer)
	}

	return &tracer
}

// Trace is a shortcut for NewTracer(opts...).Trace()
func Trace(opts ...TracerOption) func(http.Handler) http.Handler {
	return NewTracer(opts...).Trace()
}

// Trace traces every request that goes through.
func (t Tracer) Trace(opts ...TracerOption) func(http.Handler) http.Handler {
	for _, opt := range opts {
		opt(&t)
	}

	const (
		attributeServerMethod     = "http.method"
		attributeServerRoute      = "http.route"
		attributeServerStatusCode = "http.status_code"
		attributeServerPath       = "http.path"
		attributeServerUserAgent  = "http.user_agent"
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			const spanName = "delivery.http"
			var span *trace.Span

			ctx := r.Context()

			startOpts := []trace.StartOption{
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithSampler(t.sampler),
			}

			if spanCtx, ok := t.propagation.SpanContextFromRequest(r); ok {
				ctx, span = trace.StartSpanWithRemoteParent(ctx, spanName, spanCtx, startOpts...)
			} else {
				ctx, span = trace.StartSpan(ctx, spanName, startOpts...)
			}
			defer span.End()

			r = r.WithContext(ctx)
			next.ServeHTTP(rw, r)

			span.AddAttributes(
				trace.StringAttribute(attributeServerMethod, r.Method),
				trace.StringAttribute(attributeServerPath, r.URL.Path),
				trace.StringAttribute(attributeServerUserAgent, r.UserAgent()),
				trace.StringAttribute(attributeServerRoute, httpinfo.Route(r)),
				trace.Int64Attribute(attributeServerStatusCode, int64(httpinfo.Status(r))),
			)

			for _, fct := range t.callback {
				fct(r, span)
			}
		})
	}
}
