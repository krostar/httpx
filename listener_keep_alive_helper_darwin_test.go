// +build darwin

package httpx

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func getKeepAliveConfig(t *testing.T, fd int) (bool, int) {
	activated, err := syscall.GetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_KEEPALIVE)
	require.NoError(t, err)

	secs, err := syscall.GetsockoptInt(fd, syscall.IPPROTO_TCP, 0x101)
	require.NoError(t, err)

	return activated > 0, secs
}
