package rcd

import (
	"fmt"
	"os"
	"path/filepath"
)

type statFunc func(string) (os.FileInfo, error)

func findServiceScript(service string, dirs []string, stat statFunc) (string, error) {
	for _, dir := range dirs {
		p := filepath.Join(dir, service)
		if _, err := stat(p); err == nil {
			return p, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
	}
	return "", fmt.Errorf("%w: %s", ErrServiceNotFound, service)
}
