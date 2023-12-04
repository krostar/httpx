//go:build aix || freebsd || linux || netbsd
// +build aix freebsd linux netbsd

package httpx

import (
	"syscall"
	"testing"

	"gotest.tools/v3/assert"
)

func getKeepAliveConfig(t *testing.T, fd int) (bool, int) {
	activated, err := syscall.GetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_KEEPALIVE)
	assert.Check(t, err == nil)

	secs, err := syscall.GetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPINTVL)
	assert.Check(t, err == nil)

	return activated > 0, secs
}
