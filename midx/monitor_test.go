package midx

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/krostar/httpinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

func TestMonitorViewsAreRegistered(t *testing.T) {
	// reset views
	view.Unregister(httpViews...)

	var registeredViews = map[string]struct {
		mesure stats.Measure
		aggr   *view.Aggregation
		tags   []tag.Key
	}{
		"delivery/http/hit": {
			mesure: metricHTTPHit,
			aggr:   view.Count(),
			tags:   []tag.Key{tagKeyHTTPStatus, tagKeyHTTPRoute},
		},
		"delivery/http/errors": {
			mesure: metricHTTPErrors,
			aggr:   view.Count(),
			tags:   []tag.Key{tagKeyHTTPStatus, tagKeyHTTPRoute},
		},
		"delivery/http/latency": {
			mesure: metricHTTPLatency,
			aggr:   distributionHTTPLatency,
			tags:   []tag.Key{tagKeyHTTPStatus, tagKeyHTTPRoute},
		},
	}

	Monitor()

	for name, expect := range registeredViews {
		var (
			name   = name
			expect = expect
		)
		t.Run(name, func(t *testing.T) {
			v := view.Find(name)
			require.NotNil(t, v)

			assert.Equal(t, expect.mesure, v.Measure)
			assert.Equal(t, expect.aggr.Type, v.Aggregation.Type)
			assert.ElementsMatch(t, expect.tags, v.TagKeys)
		})
	}
}

func TestMonitor_no_error(t *testing.T) {
	// reset views
	view.Unregister(httpViews...)

	httpinfo.Record()(
		Monitor()(
			http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				time.Sleep(15 * time.Millisecond) // choose latency bucket
				rw.WriteHeader(http.StatusAccepted)
			}),
		),
	).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/yolo", nil))

	expectedTags := []tag.Tag{
		{Key: tagKeyHTTPRoute, Value: "GET /yolo"},
		{Key: tagKeyHTTPStatus, Value: strconv.Itoa(http.StatusAccepted)},
	}

	rows, err := view.RetrieveData("delivery/http/latency")
	require.NoError(t, err)
	assert.Len(t, rows, 1)

	rows, err = view.RetrieveData("delivery/http/hit")
	require.NoError(t, err)
	assert.Equal(t, []*view.Row{{Data: &view.CountData{Value: 1}, Tags: expectedTags}}, rows)

	rows, err = view.RetrieveData("delivery/http/errors")
	require.NoError(t, err)
	assert.Empty(t, rows)
}

func TestMonitor_error(t *testing.T) {
	// reset views
	view.Unregister(httpViews...)

	httpinfo.Record()(
		Monitor(MonitorWithHasRequestFailedFunc(func(r *http.Request) bool {
			return true
		}))(
			http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				time.Sleep(15 * time.Millisecond) // choose latency bucket
				rw.WriteHeader(http.StatusAccepted)
			}),
		),
	).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/yolo", nil))

	expectedTags := []tag.Tag{
		{Key: tagKeyHTTPRoute, Value: "GET /yolo"},
		{Key: tagKeyHTTPStatus, Value: strconv.Itoa(http.StatusAccepted)},
	}

	rows, err := view.RetrieveData("delivery/http/latency")
	require.NoError(t, err)
	assert.Len(t, rows, 1)

	rows, err = view.RetrieveData("delivery/http/hit")
	require.NoError(t, err)
	assert.Equal(t, []*view.Row{{Data: &view.CountData{Value: 1}, Tags: expectedTags}}, rows)

	rows, err = view.RetrieveData("delivery/http/errors")
	require.NoError(t, err)
	assert.Equal(t, []*view.Row{{Data: &view.CountData{Value: 1}, Tags: expectedTags}}, rows)
}
