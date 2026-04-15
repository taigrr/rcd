package rcd

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestIsRCD(t *testing.T) {
	// On non-BSD platforms, IsRCD should return false.
	// On BSD platforms, it depends on the system.
	// Either way, it shouldn't panic.
	_ = IsRCD()
}

func TestOptionsServiceDir(t *testing.T) {
	tests := []struct {
		name string
		opts Options
		want string
	}{
		{
			name: "default empty",
			opts: Options{},
			want: "",
		},
		{
			name: "custom dir",
			opts: Options{ServiceDir: "/opt/rc.d"},
			want: "/opt/rc.d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opts.ServiceDir; got != tt.want {
				t.Errorf("ServiceDir = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUnitStruct(t *testing.T) {
	u := Unit{
		Name:    "nginx",
		Path:    "/usr/local/etc/rc.d/nginx",
		Enabled: true,
	}
	if u.Name != "nginx" {
		t.Errorf("Name = %q, want %q", u.Name, "nginx")
	}
	if u.Path != "/usr/local/etc/rc.d/nginx" {
		t.Errorf("Path = %q, want %q", u.Path, "/usr/local/etc/rc.d/nginx")
	}
	if !u.Enabled {
		t.Error("Enabled = false, want true")
	}
}

func TestErrorMessages(t *testing.T) {
	errs := []struct {
		err  error
		want string
	}{
		{ErrServiceNotFound, "service not found"},
		{ErrExecTimeout, "command timed out"},
		{ErrInsufficientPermissions, "insufficient permissions"},
		{ErrMasked, "service masked"},
		{ErrNotInstalled, "service command not in $PATH"},
		{ErrSysrcNotFound, "sysrc not in $PATH"},
		{ErrServiceNotActive, "service not active"},
		{ErrUnspecified, "unknown error, please submit an issue at github.com/taigrr/rcd"},
	}
	for _, tt := range errs {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrorsAreDistinct(t *testing.T) {
	sentinels := []error{
		ErrServiceNotFound,
		ErrExecTimeout,
		ErrInsufficientPermissions,
		ErrMasked,
		ErrNotInstalled,
		ErrSysrcNotFound,
		ErrServiceNotActive,
		ErrUnspecified,
	}
	for i, a := range sentinels {
		for j, b := range sentinels {
			if i != j && errors.Is(a, b) {
				t.Errorf("expected distinct errors, but %q == %q", a, b)
			}
		}
	}
}

func TestStartNoOp(t *testing.T) {
	ctx := context.Background()
	if err := Start(ctx, "nginx", Options{}); err != nil {
		t.Errorf("Start() = %v, want nil on non-BSD", err)
	}
}

func TestStopNoOp(t *testing.T) {
	ctx := context.Background()
	if err := Stop(ctx, "nginx", Options{}); err != nil {
		t.Errorf("Stop() = %v, want nil on non-BSD", err)
	}
}

func TestRestartNoOp(t *testing.T) {
	ctx := context.Background()
	if err := Restart(ctx, "nginx", Options{}); err != nil {
		t.Errorf("Restart() = %v, want nil on non-BSD", err)
	}
}

func TestStatusNoOp(t *testing.T) {
	ctx := context.Background()
	got, err := Status(ctx, "nginx", Options{})
	if err != nil {
		t.Errorf("Status() error = %v, want nil on non-BSD", err)
	}
	if got != "" {
		t.Errorf("Status() = %q, want empty string on non-BSD", got)
	}
}

func TestIsActiveNoOp(t *testing.T) {
	ctx := context.Background()
	active, err := IsActive(ctx, "nginx", Options{})
	if err != nil {
		t.Errorf("IsActive() error = %v, want nil on non-BSD", err)
	}
	if active {
		t.Error("IsActive() = true, want false on non-BSD")
	}
}

func TestEnableNoOp(t *testing.T) {
	ctx := context.Background()
	if err := Enable(ctx, "nginx", Options{}); err != nil {
		t.Errorf("Enable() = %v, want nil on non-BSD", err)
	}
}

func TestDisableNoOp(t *testing.T) {
	ctx := context.Background()
	if err := Disable(ctx, "nginx", Options{}); err != nil {
		t.Errorf("Disable() = %v, want nil on non-BSD", err)
	}
}

func TestIsEnabledNoOp(t *testing.T) {
	ctx := context.Background()
	enabled, err := IsEnabled(ctx, "nginx", Options{})
	if err != nil {
		t.Errorf("IsEnabled() error = %v, want nil on non-BSD", err)
	}
	if enabled {
		t.Error("IsEnabled() = true, want false on non-BSD")
	}
}

func TestMaskNoOp(t *testing.T) {
	ctx := context.Background()
	if err := Mask(ctx, "nginx", Options{}); err != nil {
		t.Errorf("Mask() = %v, want nil on non-BSD", err)
	}
}

func TestUnmaskNoOp(t *testing.T) {
	ctx := context.Background()
	if err := Unmask(ctx, "nginx", Options{}); err != nil {
		t.Errorf("Unmask() = %v, want nil on non-BSD", err)
	}
}

func TestIsMaskedNoOp(t *testing.T) {
	ctx := context.Background()
	masked, err := IsMasked(ctx, "nginx", Options{})
	if err != nil {
		t.Errorf("IsMasked() error = %v, want nil on non-BSD", err)
	}
	if masked {
		t.Error("IsMasked() = true, want false on non-BSD")
	}
}

func TestReloadNoOp(t *testing.T) {
	ctx := context.Background()
	if err := Reload(ctx, "nginx", Options{}); err != nil {
		t.Errorf("Reload() = %v, want nil on non-BSD", err)
	}
}

func TestRcVarNoOp(t *testing.T) {
	ctx := context.Background()
	got, err := RcVar(ctx, "nginx", Options{})
	if err != nil {
		t.Errorf("RcVar() error = %v, want nil on non-BSD", err)
	}
	if got != "" {
		t.Errorf("RcVar() = %q, want empty string on non-BSD", got)
	}
}

func TestListNoOp(t *testing.T) {
	ctx := context.Background()
	units, err := List(ctx, Options{})
	if err != nil {
		t.Errorf("List() error = %v, want nil on non-BSD", err)
	}
	if units != nil {
		t.Errorf("List() = %v, want nil on non-BSD", units)
	}
}

func TestScriptPathNoOp(t *testing.T) {
	got, err := ScriptPath("nginx", Options{})
	if err != nil {
		t.Errorf("ScriptPath() error = %v, want nil on non-BSD", err)
	}
	if got != "" {
		t.Errorf("ScriptPath() = %q, want empty string on non-BSD", got)
	}
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// On non-BSD, all functions return nil/zero regardless of context state.
	// This test ensures no panic occurs with a cancelled context.
	if err := Start(ctx, "nginx", Options{}); err != nil {
		t.Errorf("Start() with cancelled context = %v, want nil on non-BSD", err)
	}
	if err := Stop(ctx, "nginx", Options{}); err != nil {
		t.Errorf("Stop() with cancelled context = %v, want nil on non-BSD", err)
	}
}

func TestContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(time.Millisecond) // let it expire

	// On non-BSD, should still be a no-op.
	if err := Restart(ctx, "nginx", Options{}); err != nil {
		t.Errorf("Restart() with expired context = %v, want nil on non-BSD", err)
	}
}

func TestCustomServiceDir(t *testing.T) {
	ctx := context.Background()
	opts := Options{ServiceDir: "/nonexistent/rc.d"}

	// On non-BSD, the custom service dir is ignored (no-op).
	if err := Start(ctx, "nginx", opts); err != nil {
		t.Errorf("Start() with custom ServiceDir = %v, want nil on non-BSD", err)
	}
	got, err := ScriptPath("nginx", opts)
	if err != nil {
		t.Errorf("ScriptPath() with custom ServiceDir error = %v, want nil on non-BSD", err)
	}
	if got != "" {
		t.Errorf("ScriptPath() with custom ServiceDir = %q, want empty on non-BSD", got)
	}
}

func TestEmptyServiceName(t *testing.T) {
	ctx := context.Background()
	opts := Options{}

	// On non-BSD, empty service name should still work (no-op).
	if err := Start(ctx, "", opts); err != nil {
		t.Errorf("Start() with empty name = %v, want nil on non-BSD", err)
	}
	active, err := IsActive(ctx, "", opts)
	if err != nil {
		t.Errorf("IsActive() with empty name error = %v, want nil on non-BSD", err)
	}
	if active {
		t.Error("IsActive() with empty name = true, want false on non-BSD")
	}
}
