// Copyright © 2018 Trevor N. Suarez (Rican7)

// Package version provides application version information
package version

import (
	"fmt"
	"runtime"
)

// AppName is the name of the application.
const AppName = "define"

// devID defines the default ID for development
const devID = "dev"

// This is intended to be filled by the compiler.
var (
	// ID is the VCS tag name.
	identifier = devID

	// commitHash is the VCS commit hash.
	commitHash string
)

// Name returns the name of the version.
func Name() string {
	if devID == identifier && "" != commitHash {
		return fmt.Sprintf("%s#%s", identifier, commitHash)
	}

	return identifier
}

// Printable returns a formatted printable string of the version.
func Printable() string {
	return fmt.Sprintf("%s %s (%s/%s)", AppName, Name(), runtime.GOOS, runtime.GOARCH)
}
