package httpx

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewListener(t *testing.T) {
	l, err := NewListener(context.Background(), "localhost:0")
	require.NoError(t, err)
	defer l.Close() // nolint: errcheck

	go func() {
		conn, err := net.Dial("tcp4", l.Addr().String())
		require.NoError(t, err)
		defer conn.Close() // nolint: errcheck

		_, err = io.WriteString(conn, "hello world")
		require.NoError(t, err)
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	defer conn.Close() // nolint: errcheck

	read, err := ioutil.ReadAll(conn)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(read))
}

func Test_NewListener_with_tls(t *testing.T) {
	ctx := context.Background()
	l, err := NewListener(ctx, "localhost:0", ListenWithModernTLSConfig("./testdata/cert.crt", "./testdata/cert.key"))
	require.NoError(t, err)
	defer l.Close() // nolint: errcheck

	go func() {
		rootCAPEM, err := ioutil.ReadFile("./testdata/ca.crt")
		require.NoError(t, err)
		rootCAs := x509.NewCertPool()
		require.True(t, rootCAs.AppendCertsFromPEM(rootCAPEM))

		conn, err := tls.Dial("tcp4", l.Addr().String(), &tls.Config{RootCAs: rootCAs, ServerName: "go-test"})
		require.NoError(t, err)
		defer conn.Close() // nolint: errcheck

		_, err = io.WriteString(conn, "hello world")
		require.NoError(t, err)
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	defer conn.Close() // nolint: errcheck

	read, err := ioutil.ReadAll(conn)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(read))
}

func Test_NewListener_with_bad_tls_config(t *testing.T) {
	ctx := context.Background()
	_, err := NewListener(ctx, "localhost:0", ListenWithModernTLSConfig("dont/exist", "./testdata/cert.key"))
	require.Error(t, err)
}
