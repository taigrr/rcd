package rcd

// Options configures how service commands are executed.
type Options struct {
	// ServiceDir overrides the default rc.d search paths.
	// If empty, /usr/local/etc/rc.d and /etc/rc.d are checked in order.
	ServiceDir string
}

// Unit represents a service discovered in the rc.d directories.
type Unit struct {
	// Name is the service name (e.g. "nginx", "sshd").
	Name string
	// Path is the absolute path to the rc.d script.
	Path string
	// Enabled is true if <name>_enable is set to YES in rc.conf.
	Enabled bool
}
