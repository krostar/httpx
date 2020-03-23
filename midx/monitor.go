// Package midx ...
package midx

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/krostar/httpinfo"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// nolint: gochecknoglobals
var (
	tagKeyHTTPStatus, _ = tag.NewKey("http.status")
	tagKeyHTTPRoute, _  = tag.NewKey("http.route")

	metricHTTPHit     = stats.Int64("delivery/http/hit", "Number of HTTP requests", stats.UnitDimensionless)
	metricHTTPErrors  = stats.Int64("delivery/http/errors", "Number of failing HTTP requests", stats.UnitDimensionless)
	metricHTTPLatency = stats.Float64("delivery/http/latency", "HTTP request latency", stats.UnitMilliseconds)

	distributionHTTPLatency = view.Distribution(
		1, 5, 10, 15, 30, 50, 80, 100, 300, 500, 800, 1000, 3000, 7000, 10000,
	)

	httpViews = []*view.View{{
		Name: metricHTTPHit.Name(), Description: metricHTTPHit.Description(),
		Measure:     metricHTTPHit,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{tagKeyHTTPStatus, tagKeyHTTPRoute},
	}, {
		Name: metricHTTPErrors.Name(), Description: metricHTTPErrors.Description(),
		Measure:     metricHTTPErrors,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{tagKeyHTTPStatus, tagKeyHTTPRoute},
	}, {
		Name: metricHTTPLatency.Name(), Description: metricHTTPLatency.Description(),
		Measure:     metricHTTPLatency,
		Aggregation: distributionHTTPLatency,
		TagKeys:     []tag.Key{tagKeyHTTPStatus, tagKeyHTTPRoute},
	}}
)

// Monitor monitors HTTP requests. Each requests goes through this
// middleware which records basic information about it (latency, error, ...)
func Monitor(opts ...MonitorOption) func(http.Handler) http.Handler {
	if v := view.Find(httpViews[0].Name); v == nil {
		if err := view.Register(httpViews...); err != nil {
			panic(fmt.Errorf("unable to register http views: %w", err))
		}
	}

	var o = monitorOptions{
		hasRequestFailed: func(r *http.Request) bool { return false },
	}

	for _, opt := range opts {
		opt(&o)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			latency := float64(time.Since(start).Nanoseconds()) / float64(time.Millisecond)

			ctx := r.Context()
			ctx, _ = tag.New(ctx, // nolint: errcheck
				tag.Upsert(tagKeyHTTPRoute, httpinfo.Route(r)),
				tag.Upsert(tagKeyHTTPStatus, strconv.Itoa(httpinfo.Status(r))),
			)

			stats.Record(ctx,
				metricHTTPHit.M(1),
				metricHTTPLatency.M(latency),
			)

			if o.hasRequestFailed(r) {
				stats.Record(ctx, metricHTTPErrors.M(1))
			}
		})
	}
}
