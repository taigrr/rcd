package rcd

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
)

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
		return stdout.String(), stderr.String(), code, errors.Join(ErrExecTimeout, ctx.Err())
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
