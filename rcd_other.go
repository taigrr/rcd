//go:build !freebsd && !netbsd && !openbsd && !dragonfly

package rcd

import (
	"context"
)

func start(_ context.Context, _ string, _ Options) error {
	return nil
}

func stop(_ context.Context, _ string, _ Options) error {
	return nil
}

func restart(_ context.Context, _ string, _ Options) error {
	return nil
}

func status(_ context.Context, _ string, _ Options) (string, error) {
	return "", nil
}

func isActive(_ context.Context, _ string, _ Options) (bool, error) {
	return false, nil
}

func enable(_ context.Context, _ string, _ Options) error {
	return nil
}

func disable(_ context.Context, _ string, _ Options) error {
	return nil
}

func isEnabled(_ context.Context, _ string, _ Options) (bool, error) {
	return false, nil
}

func mask(_ context.Context, _ string, _ Options) error {
	return nil
}

func unmask(_ context.Context, _ string, _ Options) error {
	return nil
}

func isMasked(_ string, _ Options) (bool, error) {
	return false, nil
}

func reload(_ context.Context, _ string, _ Options) error {
	return nil
}

func rcvar(_ context.Context, _ string, _ Options) (string, error) {
	return "", nil
}

func list(_ context.Context, _ Options) ([]Unit, error) {
	return nil, nil
}

func scriptPath(_ string, _ Options) (string, error) {
	return "", nil
}

func isRCD() bool {
	return false
}
