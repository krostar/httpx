package midx

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonitorWithHasRequestFailedFunc(t *testing.T) {
	var o monitorOptions
	MonitorWithHasRequestFailedFunc(func(*http.Request) bool { return true })(&o)
	assert.True(t, o.hasRequestFailed(nil))
}
