package httpx

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTCPListenerKeepAlive_Accept(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close() // nolint: errcheck

	lka := tcpListenerKeepAlive{
		TCPListener: l.(*net.TCPListener),
		period:      time.Second * 7,
	}

	t.Run("without keepalive", func(t *testing.T) {
		activated, _ := tcpGetKeepAliveSockOPT(t, l)
		assert.False(t, activated, "keepalive should not be set with normal listener")
	})

	t.Run("with keepalive", func(t *testing.T) {
		activated, period := tcpGetKeepAliveSockOPT(t, lka)
		assert.True(t, activated, "keepalive should have been set")
		assert.Equal(t, 7, period)
	})
}

func tcpGetKeepAliveSockOPT(t *testing.T, l net.Listener) (bool, int) {
	t.Helper()
	go func() {
		conn, err := net.Dial("tcp", l.Addr().String())
		require.NoError(t, err)
		defer conn.Close() // nolint: errcheck
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	defer conn.Close() // nolint: errcheck

	var (
		activated bool
		secs      int
	)

	r, err := conn.(*net.TCPConn).SyscallConn()
	require.NoError(t, err)
	require.NoError(t, r.Control(func(fd uintptr) {
		activated, secs = getKeepAliveConfig(t, int(fd))
	}))

	return activated, secs
}
