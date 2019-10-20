package version

import "fmt"

// The git commit compiled into the final binaries.
// This will be filled by the compiler.
var GitCommit string

const (
	// Name of the CLI application
	AppName = "okta-admin"

	// Application version
	Version = "1.0.0"
)

func FormattedVersion() string {
	version := fmt.Sprintf("Okta Admin v%s", Version)
	if GitCommit != "" {
		return fmt.Sprintf("%s (%s)", version, GitCommit)
	}
	return version
}
