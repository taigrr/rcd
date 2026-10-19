package rcd

import (
	"errors"
	"strings"
)

func parseEnabledValue(stdout string, err error) (bool, error) {
	if err == nil {
		return strings.EqualFold(strings.TrimSpace(stdout), "yes"), nil
	}
	if errors.Is(err, ErrExecTimeout) || errors.Is(err, ErrInsufficientPermissions) {
		return false, err
	}
	return false, nil
}
