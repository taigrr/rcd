package rcd

import (
	"context"
	"errors"
	"os"
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
		{ErrInvalidServiceName, "invalid service name"},
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
		ErrInvalidServiceName,
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

func TestValidateServiceName(t *testing.T) {
	tests := []struct {
		name    string
		service string
		wantErr bool
	}{
		{name: "simple", service: "sshd"},
		{name: "hyphen", service: "nginx-worker"},
		{name: "underscore", service: "postgresql_enable"},
		{name: "empty", service: "", wantErr: true},
		{name: "dot", service: ".", wantErr: true},
		{name: "parent", service: "..", wantErr: true},
		{name: "relative path", service: "../sshd", wantErr: true},
		{name: "subdirectory", service: "local/sshd", wantErr: true},
		{name: "backslash", service: `local\sshd`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateServiceName(tt.service)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidServiceName) {
					t.Fatalf("validateServiceName(%q) error = %v, want ErrInvalidServiceName", tt.service, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("validateServiceName(%q) error = %v, want nil", tt.service, err)
			}
		})
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

	if err := Start(ctx, "", opts); !errors.Is(err, ErrInvalidServiceName) {
		t.Errorf("Start() with empty name = %v, want ErrInvalidServiceName", err)
	}
	active, err := IsActive(ctx, "", opts)
	if !errors.Is(err, ErrInvalidServiceName) {
		t.Errorf("IsActive() with empty name error = %v, want ErrInvalidServiceName", err)
	}
	if active {
		t.Error("IsActive() with empty name = true, want false")
	}
}

func TestInvalidServiceNamePublicAPIs(t *testing.T) {
	ctx := context.Background()
	opts := Options{}

	errTests := []struct {
		name string
		call func() error
	}{
		{name: "Start", call: func() error { return Start(ctx, "../sshd", opts) }},
		{name: "Stop", call: func() error { return Stop(ctx, "../sshd", opts) }},
		{name: "Restart", call: func() error { return Restart(ctx, "../sshd", opts) }},
		{name: "Enable", call: func() error { return Enable(ctx, "../sshd", opts) }},
		{name: "Disable", call: func() error { return Disable(ctx, "../sshd", opts) }},
		{name: "Mask", call: func() error { return Mask(ctx, "../sshd", opts) }},
		{name: "Unmask", call: func() error { return Unmask(ctx, "../sshd", opts) }},
		{name: "Reload", call: func() error { return Reload(ctx, "../sshd", opts) }},
	}
	for _, tt := range errTests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); !errors.Is(err, ErrInvalidServiceName) {
				t.Fatalf("%s() error = %v, want ErrInvalidServiceName", tt.name, err)
			}
		})
	}

	valueTests := []struct {
		name string
		call func() error
	}{
		{name: "Status", call: func() error { _, err := Status(ctx, "../sshd", opts); return err }},
		{name: "IsActive", call: func() error { _, err := IsActive(ctx, "../sshd", opts); return err }},
		{name: "IsEnabled", call: func() error { _, err := IsEnabled(ctx, "../sshd", opts); return err }},
		{name: "IsMasked", call: func() error { _, err := IsMasked(ctx, "../sshd", opts); return err }},
		{name: "RcVar", call: func() error { _, err := RcVar(ctx, "../sshd", opts); return err }},
		{name: "ScriptPath", call: func() error { _, err := ScriptPath("../sshd", opts); return err }},
	}
	for _, tt := range valueTests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); !errors.Is(err, ErrInvalidServiceName) {
				t.Fatalf("%s() error = %v, want ErrInvalidServiceName", tt.name, err)
			}
		})
	}
}

func TestMaskedMode(t *testing.T) {
	tests := []struct {
		name string
		mode os.FileMode
		want os.FileMode
	}{
		{
			name: "preserves read write bits while clearing exec",
			mode: 0o754,
			want: 0o644,
		},
		{
			name: "keeps already masked file unchanged",
			mode: 0o640,
			want: 0o640,
		},
		{
			name: "preserves special bits",
			mode: os.ModeSetuid | 0o755,
			want: os.ModeSetuid | 0o644,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maskedMode(tt.mode); got != tt.want {
				t.Fatalf("maskedMode(%#o) = %#o, want %#o", tt.mode, got, tt.want)
			}
		})
	}
}

func TestUnmaskedMode(t *testing.T) {
	tests := []struct {
		name string
		mode os.FileMode
		want os.FileMode
	}{
		{
			name: "restores exec bits from read bits",
			mode: 0o644,
			want: 0o755,
		},
		{
			name: "keeps missing group and other read bits non-executable",
			mode: 0o600,
			want: 0o700,
		},
		{
			name: "preserves special bits",
			mode: os.ModeSetgid | 0o640,
			want: os.ModeSetgid | 0o750,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unmaskedMode(tt.mode); got != tt.want {
				t.Fatalf("unmaskedMode(%#o) = %#o, want %#o", tt.mode, got, tt.want)
			}
		})
	}
}
