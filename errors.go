package rcd

import (
	"errors"
)

var (
	// ErrServiceNotFound is returned when the rc.d script does not exist
	// in either /usr/local/etc/rc.d or /etc/rc.d.
	ErrServiceNotFound = errors.New("service not found")

	// ErrExecTimeout is returned when the provided context was cancelled
	// before the command finished execution.
	ErrExecTimeout = errors.New("command timed out")

	// ErrInsufficientPermissions is returned when the command requires
	// elevated privileges (typically root).
	ErrInsufficientPermissions = errors.New("insufficient permissions")

	// ErrMasked is returned when an operation is attempted on a masked
	// service. A service is considered masked when its rc.d script has
	// had its execute permission removed.
	ErrMasked = errors.New("service masked")

	// ErrNotInstalled is returned when the `service` command cannot be
	// found in $PATH.
	ErrNotInstalled = errors.New("service command not in $PATH")

	// ErrSysrcNotFound is returned when sysrc(8) is not available.
	// sysrc is required for enable/disable operations.
	ErrSysrcNotFound = errors.New("sysrc not in $PATH")

	// ErrServiceNotActive is returned when a service was expected to be
	// running but is not.
	ErrServiceNotActive = errors.New("service not active")

	// ErrUnspecified is a catch-all for unrecognized errors. If you
	// encounter this, please submit an issue at github.com/taigrr/rcd.
	ErrUnspecified = errors.New("unknown error, please submit an issue at github.com/taigrr/rcd")
)
