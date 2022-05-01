package server

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidEOFByte = errors.New("invalid end of file byte")
)
