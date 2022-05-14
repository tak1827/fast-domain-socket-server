package server

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tak1827/fast-domain-socket-server/data"
	"go.uber.org/goleak"
)

func client(tx data.Message) (err error) {
	var (
		dst = make([]byte, 1024)
	)
	conn, err := net.Dial("unix", DefaultSockFilePath)
	if err != nil {
		return
	}
	defer conn.Close()

	d, _ := tx.Marshal()
	d = append(d, EOFByte)
	if _, err = conn.Write(d); err != nil {
		return
	}

	if dst, err = ReadConn(conn, dst); err != nil {
		return
	}

	if err = tx.Unmarshal(dst); err != nil {
		return
	}

	return
}

func TestShutdown(t *testing.T) {
	defer goleak.VerifyNone(t)

	s := NewServer(DefaultSockFilePath)
	ln, err := s.Listen()
	require.NoError(t, err)

	go func() {
		require.NoError(t, s.Serve(ln))
	}()

	time.Sleep(1 * time.Second)
	require.NoError(t, s.Shutdown(ln))
}

func TestClient(t *testing.T) {
	defer goleak.VerifyNone(t)

	var (
		s  = NewServer(DefaultSockFilePath)
		tx = data.Message{
			Type:    "type",
			Payload: "palyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyad",
		}
	)
	ln, err := s.Listen()
	require.NoError(t, err)
	defer s.Shutdown(ln)

	go func() {
		require.NoError(t, s.Serve(ln))
	}()

	for n := 0; n < 100; n++ {
		go func() {
			client(tx)
		}()
	}
}

func BenchmarkServer(b *testing.B) {
	var (
		s  = NewServer(DefaultSockFilePath)
		tx = data.Message{
			Type:    "type",
			Payload: "palyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyadpalyad",
		}
	)
	ln, err := s.Listen()
	require.NoError(b, err)
	defer s.Shutdown(ln)

	go func() {
		require.NoError(b, s.Serve(ln))
	}()

	for n := 0; n < b.N; n++ {
		if err = client(tx); err != nil {
			b.Fatal(err)
		}
	}
}
