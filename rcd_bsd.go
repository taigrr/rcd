//go:build freebsd || netbsd || openbsd || dragonfly

package rcd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	serviceBin string
	sysrcBin   string
)

func init() {
	serviceBin, _ = exec.LookPath("service")
	sysrcBin, _ = exec.LookPath("sysrc")
}

// rcDirs returns the directories to search for rc.d scripts in priority order.
func rcDirs(opts Options) []string {
	if opts.ServiceDir != "" {
		return []string{opts.ServiceDir}
	}
	return []string{"/usr/local/etc/rc.d", "/etc/rc.d"}
}

func scriptPath(service string, opts Options) (string, error) {
	for _, dir := range rcDirs(opts) {
		p := filepath.Join(dir, service)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrServiceNotFound, service)
}

func execute(ctx context.Context, name string, args ...string) (string, string, int, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if cmd.ProcessState != nil {
		code = cmd.ProcessState.ExitCode()
	}

	if ctx.Err() != nil {
		return stdout.String(), stderr.String(), code, ErrExecTimeout
	}

	if err != nil {
		stderrStr := stderr.String()
		customErr := filterErr(stderrStr)
		if customErr != nil {
			return stdout.String(), stderrStr, code, customErr
		}
		return stdout.String(), stderrStr, code, err
	}

	return stdout.String(), stderr.String(), code, nil
}

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

func serviceCmd(ctx context.Context, service, action string) (string, string, int, error) {
	if serviceBin == "" {
		return "", "", 1, ErrNotInstalled
	}
	return execute(ctx, serviceBin, service, action)
}

func start(ctx context.Context, service string, opts Options) error {
	if _, err := scriptPath(service, opts); err != nil {
		return err
	}
	_, _, _, err := serviceCmd(ctx, service, "start")
	return err
}

func stop(ctx context.Context, service string, opts Options) error {
	if _, err := scriptPath(service, opts); err != nil {
		return err
	}
	_, _, _, err := serviceCmd(ctx, service, "stop")
	return err
}

func restart(ctx context.Context, service string, opts Options) error {
	if _, err := scriptPath(service, opts); err != nil {
		return err
	}
	_, _, _, err := serviceCmd(ctx, service, "restart")
	return err
}

func status(ctx context.Context, service string, opts Options) (string, error) {
	if _, err := scriptPath(service, opts); err != nil {
		return "", err
	}
	stdout, _, _, err := serviceCmd(ctx, service, "status")
	return strings.TrimSpace(stdout), err
}

func isActive(ctx context.Context, service string, opts Options) (bool, error) {
	if _, err := scriptPath(service, opts); err != nil {
		return false, err
	}
	_, _, code, err := serviceCmd(ctx, service, "status")
	if code == 1 {
		// Exit code 1 = not running.
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return code == 0, nil
}

func rcVarName(service string) string {
	return service + "_enable"
}

func enable(ctx context.Context, service string, opts Options) error {
	if _, err := scriptPath(service, opts); err != nil {
		return err
	}
	if sysrcBin == "" {
		return ErrSysrcNotFound
	}
	_, _, _, err := execute(ctx, sysrcBin, rcVarName(service)+"=YES")
	return err
}

func disable(ctx context.Context, service string, opts Options) error {
	if _, err := scriptPath(service, opts); err != nil {
		return err
	}
	if sysrcBin == "" {
		return ErrSysrcNotFound
	}
	_, _, _, err := execute(ctx, sysrcBin, rcVarName(service)+"=NO")
	return err
}

func isEnabled(ctx context.Context, service string, opts Options) (bool, error) {
	if _, err := scriptPath(service, opts); err != nil {
		return false, err
	}
	if sysrcBin == "" {
		return false, ErrSysrcNotFound
	}
	stdout, _, _, err := execute(ctx, sysrcBin, "-n", rcVarName(service))
	if err != nil {
		// Variable not set means not enabled.
		return false, nil
	}
	return strings.EqualFold(strings.TrimSpace(stdout), "yes"), nil
}

func mask(_ context.Context, service string, opts Options) error {
	p, err := scriptPath(service, opts)
	if err != nil {
		return err
	}
	return os.Chmod(p, 0o444)
}

func unmask(_ context.Context, service string, opts Options) error {
	p, err := scriptPath(service, opts)
	if err != nil {
		return err
	}
	return os.Chmod(p, 0o755)
}

func isMasked(service string, opts Options) (bool, error) {
	p, err := scriptPath(service, opts)
	if err != nil {
		return false, err
	}
	info, err := os.Stat(p)
	if err != nil {
		return false, err
	}
	return info.Mode()&0o111 == 0, nil
}

func reload(ctx context.Context, service string, opts Options) error {
	if _, err := scriptPath(service, opts); err != nil {
		return err
	}
	_, _, _, err := serviceCmd(ctx, service, "reload")
	return err
}

func rcvar(ctx context.Context, service string, opts Options) (string, error) {
	if _, err := scriptPath(service, opts); err != nil {
		return "", err
	}
	stdout, _, _, err := serviceCmd(ctx, service, "rcvar")
	return strings.TrimSpace(stdout), err
}

func list(ctx context.Context, opts Options) ([]Unit, error) {
	var units []Unit
	seen := make(map[string]bool)

	for _, dir := range rcDirs(opts) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		for _, entry := range entries {
			name := entry.Name()
			if seen[name] || strings.HasPrefix(name, ".") {
				continue
			}
			if !entry.Type().IsRegular() && entry.Type()&os.ModeSymlink == 0 {
				continue
			}
			seen[name] = true

			enabled, _ := isEnabled(ctx, name, opts)
			units = append(units, Unit{
				Name:    name,
				Path:    filepath.Join(dir, name),
				Enabled: enabled,
			})
		}
	}
	return units, nil
}

func isRCD() bool {
	info, err := os.Stat("/etc/rc.d")
	return err == nil && info.IsDir()
}
