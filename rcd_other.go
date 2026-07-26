//go:build !freebsd && !netbsd && !openbsd && !dragonfly

package rcd

import (
	"context"
)

func start(_ context.Context, service string, _ Options) error {
	return validateServiceName(service)
}

func stop(_ context.Context, service string, _ Options) error {
	return validateServiceName(service)
}

func restart(_ context.Context, service string, _ Options) error {
	return validateServiceName(service)
}

func status(_ context.Context, service string, _ Options) (string, error) {
	return "", validateServiceName(service)
}

func isActive(_ context.Context, service string, _ Options) (bool, error) {
	return false, validateServiceName(service)
}

func enable(_ context.Context, service string, _ Options) error {
	return validateServiceName(service)
}

func disable(_ context.Context, service string, _ Options) error {
	return validateServiceName(service)
}

func isEnabled(_ context.Context, service string, _ Options) (bool, error) {
	return false, validateServiceName(service)
}

func mask(_ context.Context, service string, _ Options) error {
	return validateServiceName(service)
}

func unmask(_ context.Context, service string, _ Options) error {
	return validateServiceName(service)
}

func isMasked(service string, _ Options) (bool, error) {
	return false, validateServiceName(service)
}

func reload(_ context.Context, service string, _ Options) error {
	return validateServiceName(service)
}

func rcvar(_ context.Context, service string, _ Options) (string, error) {
	return "", validateServiceName(service)
}

func list(_ context.Context, _ Options) ([]Unit, error) {
	return nil, nil
}

func scriptPath(service string, _ Options) (string, error) {
	return "", validateServiceName(service)
}

func isRCD() bool {
	return false
}
