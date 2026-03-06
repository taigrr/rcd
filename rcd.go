// Package rcd provides idiomatic Go bindings for managing BSD rc.d services.
//
// It wraps the service(8) and sysrc(8) commands to provide a clean API for
// starting, stopping, enabling, and querying services on FreeBSD, NetBSD,
// OpenBSD, and DragonFlyBSD.
//
// On non-BSD platforms, all functions are no-ops that return nil/zero values,
// allowing cross-platform code to compile without build tags.
package rcd

import (
	"context"
)

// Start activates the specified service.
//
// Equivalent to `service <name> start`.
func Start(ctx context.Context, service string, opts Options) error {
	return start(ctx, service, opts)
}

// Stop deactivates the specified service.
//
// Equivalent to `service <name> stop`.
func Stop(ctx context.Context, service string, opts Options) error {
	return stop(ctx, service, opts)
}

// Restart stops and then starts the specified service.
//
// Equivalent to `service <name> restart`.
func Restart(ctx context.Context, service string, opts Options) error {
	return restart(ctx, service, opts)
}

// Status returns the raw output of `service <name> status`.
func Status(ctx context.Context, service string, opts Options) (string, error) {
	return status(ctx, service, opts)
}

// IsActive returns true if the service is currently running.
//
// It checks the exit code of `service <name> status`.
func IsActive(ctx context.Context, service string, opts Options) (bool, error) {
	return isActive(ctx, service, opts)
}

// Enable sets <service>_enable="YES" in rc.conf via sysrc(8).
func Enable(ctx context.Context, service string, opts Options) error {
	return enable(ctx, service, opts)
}

// Disable sets <service>_enable="NO" in rc.conf via sysrc(8).
func Disable(ctx context.Context, service string, opts Options) error {
	return disable(ctx, service, opts)
}

// IsEnabled returns true if <service>_enable is set to "YES" in rc.conf.
func IsEnabled(ctx context.Context, service string, opts Options) (bool, error) {
	return isEnabled(ctx, service, opts)
}

// Mask prevents a service from being started by removing the execute
// permission from its rc.d script. BSD does not have a formal mask
// concept like systemd; this is the closest equivalent.
func Mask(ctx context.Context, service string, opts Options) error {
	return mask(ctx, service, opts)
}

// Unmask restores the execute permission on the service's rc.d script.
func Unmask(ctx context.Context, service string, opts Options) error {
	return unmask(ctx, service, opts)
}

// IsMasked returns true if the service's rc.d script lacks execute permission.
func IsMasked(_ context.Context, service string, opts Options) (bool, error) {
	return isMasked(service, opts)
}

// Reload sends a reload signal to the specified service, if supported.
//
// Equivalent to `service <name> reload`.
func Reload(ctx context.Context, service string, opts Options) error {
	return reload(ctx, service, opts)
}

// RcVar returns the rc.conf variables associated with a service.
//
// Equivalent to `service <name> rcvar`.
func RcVar(ctx context.Context, service string, opts Options) (string, error) {
	return rcvar(ctx, service, opts)
}

// List returns all services found in the rc.d directories.
func List(ctx context.Context, opts Options) ([]Unit, error) {
	return list(ctx, opts)
}

// ScriptPath returns the absolute path to the rc.d script for a service.
// It checks /usr/local/etc/rc.d first (ports/packages), then /etc/rc.d (base).
// Returns ErrServiceNotFound if the script doesn't exist.
func ScriptPath(service string, opts Options) (string, error) {
	return scriptPath(service, opts)
}

// IsRCD returns true if the current system uses rc.d as its init system.
// It checks for the existence of /etc/rc.d.
func IsRCD() bool {
	return isRCD()
}
