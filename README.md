[![PkgGoDev](https://pkg.go.dev/badge/github.com/taigrr/rcd)](https://pkg.go.dev/github.com/taigrr/rcd)
# rcd

This library provides idiomatic Go bindings for managing BSD rc.d services, making it easier to write system tooling that targets FreeBSD, NetBSD, OpenBSD, and DragonFlyBSD.

It wraps `service(8)` and `sysrc(8)` in a structured, type-safe API — no more shelling out and parsing output by hand.

If your system isn't running rc.d, this library will compile but all functions will be no-ops.

## What is rc.d

rc.d is the service management framework used by BSD operating systems. Services are defined as shell scripts in `/etc/rc.d` (base system) and `/usr/local/etc/rc.d` (ports/packages). Configuration is managed through `/etc/rc.conf` variables.

## Supported operations

- [x] `service <name> start`
- [x] `service <name> stop`
- [x] `service <name> restart`
- [x] `service <name> reload`
- [x] `service <name> status`
- [x] `service <name> rcvar`
- [x] `sysrc <name>_enable=YES` (enable)
- [x] `sysrc <name>_enable=NO` (disable)
- [x] Mask/unmask via execute permission

## Helper functionality

- [x] Check if a service is running (`IsActive`)
- [x] Check if a service is enabled at boot (`IsEnabled`)
- [x] Check if a service is masked (`IsMasked`)
- [x] Find the rc.d script path (`ScriptPath`)
- [x] List all available services (`List`)
- [x] Detect if the system uses rc.d (`IsRCD`)

## Useful errors

All functions return predefined error types for common failure cases:

| Error | Meaning |
|-------|---------|
| `ErrServiceNotFound` | The rc.d script doesn't exist |
| `ErrExecTimeout` | Context was cancelled before completion |
| `ErrInsufficientPermissions` | Need root/elevated privileges |
| `ErrMasked` | Service is masked (execute bit removed) |
| `ErrNotInstalled` | `service` command not in `$PATH` |
| `ErrSysrcNotFound` | `sysrc` command not in `$PATH` |
| `ErrServiceNotActive` | Service expected to be running but isn't |

## Context support

All calls support Go's `context` functionality, allowing timeouts and cancellation.

## Simple example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/taigrr/rcd"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    opts := rcd.Options{}

    // Enable and start nginx
    if err := rcd.Enable(ctx, "nginx", opts); err != nil {
        log.Fatalf("unable to enable nginx: %v", err)
    }
    if err := rcd.Start(ctx, "nginx", opts); err != nil {
        log.Fatalf("unable to start nginx: %v", err)
    }

    // Check status
    active, err := rcd.IsActive(ctx, "nginx", opts)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("nginx running:", active)

    // List all services
    units, err := rcd.List(ctx, opts)
    if err != nil {
        log.Fatal(err)
    }
    for _, u := range units {
        status := " "
        if u.Enabled {
            status = "*"
        }
        fmt.Printf("[%s] %s (%s)\n", status, u.Name, u.Path)
    }
}
```

## Cross-platform compilation

On non-BSD platforms, all functions compile and return zero values. This allows you to use the library in cross-platform projects without build tags in your own code.

## License

This project is licensed under the 0BSD License, written by [Rob Landley](https://github.com/landley).
As such, you may use this library without restriction or attribution, but please don't pass it off as your own.
Attribution, though not required, is appreciated.

By contributing, you agree all code submitted also falls under the License.

## External resources

- [FreeBSD rc.d handbook](https://docs.freebsd.org/en/articles/rc-scripting/)
- [service(8) man page](https://man.freebsd.org/cgi/man.cgi?query=service)
- [sysrc(8) man page](https://man.freebsd.org/cgi/man.cgi?query=sysrc)
