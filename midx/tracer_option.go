package midx

import (
	"net/http"

	"go.opencensus.io/trace"
)

// TracerOption applies options to the tracer.
type TracerOption func(o *Tracer)

// TracerWithAlwaysSampler sets the sampler of the tracer.
func TracerWithAlwaysSampler() TracerOption {
	return TracerWithSampler(trace.AlwaysSample())
}

// TracerWithNeverSampler sets the sampler of the tracer.
func TracerWithNeverSampler() TracerOption {
	return TracerWithSampler(trace.NeverSample())
}

// TracerWithProbabilitySampler sets the sampler of the tracer.
func TracerWithProbabilitySampler(fraction float64) TracerOption {
	return TracerWithSampler(trace.ProbabilitySampler(fraction))
}

// TracerWithSampler sets the sampler of the tracer.
func TracerWithSampler(sampler trace.Sampler) TracerOption {
	return func(o *Tracer) {
		o.sampler = sampler
	}
}

// TracerWithCallback sets function that will be called
// at the end of the request in the tracer.
func TracerWithCallback(fcts ...func(r *http.Request, span *trace.Span)) TracerOption {
	return func(o *Tracer) {
		o.callback = append(o.callback, fcts...)
	}
}
