package rcd

import (
	"testing"
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
	errors := []struct {
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
	}
	for _, tt := range errors {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}
