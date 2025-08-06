package version

import (
	"fmt"
	"runtime"
)

// These variables are set at build time using ldflags
var (
	// Version is the current version, set by ldflags during build
	Version = "dev"
	// GitCommit is the git commit hash, set by ldflags during build  
	GitCommit = "unknown"
	// BuildDate is the build date, set by ldflags during build
	BuildDate = "unknown"
)

const (
	// Name is the application name
	Name = "pyinit"
)

// Info contains build and version information
type Info struct {
	Version   string
	GitCommit string
	BuildDate string
	GoVersion string
	Platform  string
}

// GetVersion returns the current version
func GetVersion() string {
	return Version
}

// GetBuildInfo returns detailed build information
func GetBuildInfo() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a formatted version string
func (i Info) String() string {
	result := fmt.Sprintf("%s %s", Name, i.Version)
	if i.GitCommit != "unknown" {
		result += fmt.Sprintf("\nCommit: %s", i.GitCommit)
	}
	if i.BuildDate != "unknown" {
		result += fmt.Sprintf("\nBuilt: %s", i.BuildDate)
	}
	result += fmt.Sprintf("\nGo: %s\nPlatform: %s", i.GoVersion, i.Platform)
	return result
}