package server

import (
	"github.com/tak1827/fast-domain-socket-server/data"
)

const (
	DefaultSockFilePath   = "../domain.sock"
	DefaultTimeout        = int64(60) // 60 sec
	DefaultReadBufferSize = 4096
)

var (
	DefaultHandler = func(tx *data.Message) ([]byte, error) {
		return tx.Marshal()
	}
	DefaultErrHandler = func(err error) {
		panic(err)
	}
)

type Option interface {
	Apply(*Server)
}

type TimeoutOpt int64

func (t TimeoutOpt) Apply(s *Server) {
	s.timeout = int64(t)
}
func WithTimeout(t int64) TimeoutOpt {
	if t <= 0 {
		panic("Timeout should be positive")
	}
	return TimeoutOpt(t)
}

type ReadBufferSizeOpt int

func (t ReadBufferSizeOpt) Apply(s *Server) {
	s.readBufferSize = int(t)
}
func WithReadBufferSize(t int) ReadBufferSizeOpt {
	if t <= 0 {
		panic("readBufferSize should be positive")
	}
	return ReadBufferSizeOpt(t)
}

type HandlerOpt Handler

func (t HandlerOpt) Apply(s *Server) {
	s.handler = Handler(t)
}
func WithHandler(t HandlerOpt) HandlerOpt {
	return HandlerOpt(t)
}

type ErrHandlerOpt ErrHandler

func (t ErrHandlerOpt) Apply(s *Server) {
	s.errHandler = ErrHandler(t)
}
func WithErrHandler(t ErrHandlerOpt) ErrHandlerOpt {
	return ErrHandlerOpt(t)
}
