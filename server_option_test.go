package httpx

import (
	"log"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func Test_ServerWithErrorLogger(t *testing.T) {
	var o serverOptions

	ServerWithErrorLogger(log.New(os.Stderr, "http", 0))(&o)
	assert.Check(t, o.errorLogger != nil)
}
