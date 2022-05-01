package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/tak1827/fast-domain-socket-server/data"
)

const (
	SockFilePath = "../domain.sock"

	EOFByte = 0x12 // -> same as "\n"

	defaultReadBufferSize  = 4096
	defaultWriteBufferSize = 4096
)

type Server struct {
	addr           string
	timeout        time.Duration
	readBufferSize int
	ln             net.Listener
	wg             sync.WaitGroup

	ErrCh chan error
}

func NewServer(addr string) (s Server) {
	s.addr = addr
	s.timeout = 3 * time.Second
	s.readBufferSize = defaultReadBufferSize
	s.ErrCh = make(chan error)
	return
}

func (s *Server) ListenAndServeUNIX() (err error) {
	if err = removeSocketFile(s.addr); err != nil {
		return
	}
	if s.ln, err = net.Listen("unix", s.addr); err != nil {
		return
	}
	if err = os.Chmod(s.addr, 0700); err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot chmod %s", s.addr))
	}
	return s.Serve()
}

func (s *Server) Serve() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}

		if err = conn.SetDeadline(time.Now().Add(s.timeout)); err != nil {
			return err
		}

		s.wg.Add(1)

		go func() {
			defer s.wg.Done()

			if err = s.serveConn(conn); err != nil && !errors.Is(err, io.EOF) {
				s.ErrCh <- err
			}

			conn.Close()
		}()
	}
}

func (s *Server) serveConn(conn net.Conn) (err error) {
	dst := make([]byte, s.readBufferSize)
	if dst, err = ReadConn(conn, dst); err != nil {
		err = errors.Wrap(err, "faild to read conn")
		return
	}

	tx := data.Message{}
	if err = tx.Unmarshal(dst); err != nil {
		return
	}

	fmt.Printf("Recieved: %v\n", tx)

	d, _ := tx.Marshal()

	d = append(d, EOFByte)

	if _, err = conn.Write(d); err != nil {
		return
	}

	return
}

func (s *Server) Shutdown() (err error) {
	if err = s.ln.Close(); err != nil {
		return err
	}

	s.wg.Wait()

	if err = removeSocketFile(s.addr); err != nil {
		return
	}

	return
}
