package midx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/krostar/httpinfo"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
	"gopkg.in/go-playground/assert.v1"
)

type mapExporter map[string]*trace.SpanData

func (e mapExporter) ExportSpan(s *trace.SpanData) { e[s.Name] = s }

func TestTracer_Trace(t *testing.T) {
	var spans = make(mapExporter)
	trace.RegisterExporter(&spans)
	defer trace.UnregisterExporter(&spans)

	httpinfo.Record()(
		Trace()(
			http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				_, span := trace.StartSpan(r.Context(), "my-super-job")
				defer span.End()
				rw.WriteHeader(http.StatusAccepted)
			}),
		),
	).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/toto", nil))

	trace.UnregisterExporter(&spans)
	require.Len(t, spans, 2)

	root, ok := spans["delivery.http"]
	require.True(t, ok)
	child, ok := spans["my-super-job"]
	require.True(t, ok)

	assert.Equal(t, root.TraceID, child.TraceID)
	assert.Equal(t, root.SpanID, child.ParentSpanID)
	assert.Equal(t, http.MethodGet, root.Attributes["http.method"])
	assert.Equal(t, "GET /toto", root.Attributes["http.route"])
	assert.Equal(t, int64(http.StatusAccepted), root.Attributes["http.status_code"])
	assert.Equal(t, "/toto", root.Attributes["http.path"])
}

func TestHTTPTracer_Trace_override_options(t *testing.T) {
	var (
		spans   = make(mapExporter)
		tracer  = NewTracer(TracerWithAlwaysSampler())
		handler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			_, span := trace.StartSpan(r.Context(), "my-super-job")
			defer span.End()
			rw.WriteHeader(http.StatusAccepted)
		})
	)

	trace.RegisterExporter(&spans)
	httpinfo.Record()(
		tracer.Trace(TracerWithNeverSampler())(handler),
	).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/toto", nil))
	trace.UnregisterExporter(&spans)
	require.Len(t, spans, 0)

	trace.RegisterExporter(&spans)
	httpinfo.Record()(
		tracer.Trace()(handler),
	).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/toto", nil))
	trace.UnregisterExporter(&spans)
	require.Len(t, spans, 2)
}
