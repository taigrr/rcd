package rcd

import "strings"

func validateServiceName(service string) error {
	if service == "" || service == "." || service == ".." || strings.HasPrefix(service, "-") || strings.ContainsAny(service, `/\`) {
		return ErrInvalidServiceName
	}
	return nil
}
