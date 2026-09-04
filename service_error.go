package rcd

import (
	"errors"
	"fmt"
	"strings"
)

func filterErr(stderr string) error {
	lower := strings.ToLower(stderr)
	switch {
	case strings.Contains(lower, "not found"):
		return errors.Join(ErrServiceNotFound, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(lower, "does not exist"):
		return errors.Join(ErrServiceNotFound, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(lower, "permission denied"):
		return errors.Join(ErrInsufficientPermissions, fmt.Errorf("stderr: %s", stderr))
	case strings.Contains(lower, "not permitted"):
		return errors.Join(ErrInsufficientPermissions, fmt.Errorf("stderr: %s", stderr))
	default:
		return nil
	}
}
