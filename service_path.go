package rcd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type statFunc func(string) (os.FileInfo, error)

func findServiceScript(service string, dirs []string, stat statFunc) (string, error) {
	var firstErr error
	for _, dir := range dirs {
		p := filepath.Join(dir, service)
		if _, err := stat(p); err == nil {
			return p, nil
		} else if !errors.Is(err, fs.ErrNotExist) && firstErr == nil {
			firstErr = err
		}
	}
	if firstErr != nil {
		return "", firstErr
	}
	return "", fmt.Errorf("%w: %s", ErrServiceNotFound, service)
}
