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
	EOFByte = 0x12 // -> same as "\n"
)

var (
	ErrInvalidEOFByte = errors.New("invalid end of file byte")

	timeoutDuration time.Duration
)

type Handler func(tx *data.Message) ([]byte, error)
type ErrHandler func(err error)

type Server struct {
	sync.Mutex

	addr           string
	timeout        int64
	readBufferSize int

	wg    sync.WaitGroup
	pool  sync.Pool
	errCh chan error

	handler    Handler
	errHandler ErrHandler
}

func NewServer(addr string, opts ...Option) (s Server) {
	s.addr = addr
	s.timeout = DefaultTimeout
	s.readBufferSize = DefaultReadBufferSize
	s.errCh = make(chan error)
	s.handler = DefaultHandler
	s.errHandler = DefaultErrHandler

	if s.addr == "" {
		s.addr = DefaultSockFilePath
	}

	for i := range opts {
		opts[i].Apply(&s)
	}

	timeoutDuration = time.Duration(time.Duration(s.timeout) * time.Second)

	return
}

func (s *Server) Listen() (net.Listener, error) {
	err := removeSocketFile(s.addr)
	if err != nil {
		return nil, err
	}

	ln, err := net.Listen("unix", s.addr)
	if err != nil {
		return nil, err
	}

	if err = os.Chmod(s.addr, 0700); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("cannot chmod %s", s.addr))
	}

	return ln, nil
}

func (s *Server) Serve(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			if isClosedConnError(err) {
				return nil
			}

			return err
		}

		if err = conn.SetDeadline(time.Now().Add(timeoutDuration)); err != nil {
			return errors.Wrap(err, "failed to set deadline")
		}

		s.Lock()
		s.wg.Add(1)
		s.Unlock()

		go func() {
			defer s.wg.Done()

			if err = s.serveConn(conn); err != nil && !errors.Is(err, io.EOF) {
				s.errCh <- err
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

	// tx := &data.Message{}
	v := s.pool.Get()
	if v == nil {
		v = &data.Message{}
	}
	tx := v.(*data.Message)
	defer s.pool.Put(tx)

	if err = tx.Unmarshal(dst); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to unmarshal packet(=%v)", dst))
	}

	data, err := s.handler(tx)
	if err != nil {
		return errors.Wrap(err, "failed to handle")
	}

	return WriteConn(conn, data)
}

func (s *Server) Shutdown(ln net.Listener) (err error) {
	if err = ln.Close(); err != nil {
		return err
	}

	s.wg.Wait()

	if err = removeSocketFile(s.addr); err != nil {
		return
	}

	close(s.errCh)

	return
}

func removeSocketFile(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, fmt.Sprintf("unexpected error when trying to remove unix socket file %q", path))
	}
	return nil
}
