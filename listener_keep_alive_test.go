package httpx

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewListener_with_keepalive(t *testing.T) {
	t.Run("without keepalive", func(t *testing.T) {
		l, err := NewListener(context.Background(), "localhost:0", ListenWithoutKeepAlive())
		require.NoError(t, err)
		defer l.Close() // nolint: errcheck

		activated, _ := tcpGetKeepAliveSockOPT(t, l)
		assert.False(t, activated, "keepalive should not be set with normal listener")
	})

	t.Run("with keepalive", func(t *testing.T) {
		l, err := NewListener(context.Background(), "localhost:0", ListenWithKeepAlive(17*time.Second))
		require.NoError(t, err)
		defer l.Close() // nolint: errcheck

		activated, period := tcpGetKeepAliveSockOPT(t, l)
		assert.True(t, activated, "keepalive should have been set")
		assert.Equal(t, 17, period)
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
