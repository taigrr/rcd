package rcd

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
)

func TestExecuteContextErrors(t *testing.T) {
	for _, deadline := range []bool{false, true} {
		name := "canceled"
		if deadline {
			name = "deadline"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			want := context.Canceled
			if deadline {
				cancel()
				ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
				want = context.DeadlineExceeded
			} else {
				cancel()
			}
			defer cancel()
			_, _, _, err := execute(ctx, os.Args[0], "-test.run=^$")
			if !errors.Is(err, ErrExecTimeout) || !errors.Is(err, want) {
				t.Fatalf("execute() error = %v, want ErrExecTimeout and %v", err, want)
			}
		})
	}
}

func TestExecuteSuccess(t *testing.T) {
	_, _, code, err := execute(context.Background(), os.Args[0], "-test.run=^$")
	if err != nil || code != 0 {
		t.Fatalf("execute() = code %d, error %v", code, err)
	}
}

func TestExecuteRunningDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	_, _, _, err := execute(ctx, os.Args[0], "-test.run=^TestExecuteSleepingHelper$", "--", "rcd-sleep")
	if !errors.Is(err, ErrExecTimeout) || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("execute() error = %v, want timeout and deadline exceeded", err)
	}
}

func TestExecuteSleepingHelper(t *testing.T) {
	if os.Args[len(os.Args)-1] != "rcd-sleep" {
		return
	}
	time.Sleep(10 * time.Second)
}
