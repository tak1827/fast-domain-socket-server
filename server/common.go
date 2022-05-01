package server

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func removeSocketFile(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, fmt.Sprintf("unexpected error when trying to remove unix socket file %q", path))
	}
	return nil
}
