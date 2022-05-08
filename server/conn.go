package server

import (
	"io"
	"net"
	"strings"

	"github.com/lithdew/bytesutil"
	"github.com/pkg/errors"
)

func ReadConn(conn net.Conn, dst []byte) ([]byte, error) {
	var (
		buf  []byte
		size int
	)
	for {
		n, err := conn.Read(dst)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		size += n

		if n == len(dst) || len(buf) != 0 {
			buf = bytesutil.ExtendSlice(buf, size)
			copy(buf[size-n:], dst[:n])

			if n == len(dst) && dst[n-1] != EOFByte {
				continue
			}

		}

		if dst[n-1] != EOFByte {
			return nil, ErrInvalidEOFByte
		}

		break
	}

	size += 1 // increment one as EOFbyte

	if len(buf) != 0 {
		return buf[:size], nil
	}

	return dst[:size], nil
}

func isClosedConnError(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}
