package httpx

import (
	"context"
	"net"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func Test_NewListener_with_keepalive(t *testing.T) {
	t.Run("without keepalive", func(t *testing.T) {
		l, err := NewListener(context.Background(), "localhost:0", ListenWithoutKeepAlive())
		assert.NilError(t, err)

		activated, _ := tcpGetKeepAliveSockOPT(t, l)
		assert.Check(t, !activated, "keepalive should not be set with normal listener")

		assert.NilError(t, l.Close())
	})

	t.Run("with keepalive", func(t *testing.T) {
		l, err := NewListener(context.Background(), "localhost:0", ListenWithKeepAlive(17*time.Second))
		assert.NilError(t, err)

		activated, period := tcpGetKeepAliveSockOPT(t, l)
		assert.Check(t, activated, "keepalive should have been set")
		assert.Equal(t, 17, period)

		assert.NilError(t, l.Close())
	})
}

func tcpGetKeepAliveSockOPT(t *testing.T, l net.Listener) (bool, int) {
	t.Helper()
	go func() {
		conn, err := net.Dial("tcp", l.Addr().String())
		assert.NilError(t, err)
		assert.NilError(t, conn.Close())
	}()

	conn, err := l.Accept()
	assert.NilError(t, err)

	var (
		activated bool
		secs      int
	)

	r, err := conn.(*net.TCPConn).SyscallConn()
	assert.NilError(t, err)
	assert.NilError(t, r.Control(func(fd uintptr) {
		activated, secs = getKeepAliveConfig(t, int(fd))
	}))

	assert.NilError(t, conn.Close())

	return activated, secs
}
