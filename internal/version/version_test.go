package version

import (
	"runtime"
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	
	// Version should not be empty
	if version == "" {
		t.Error("GetVersion() returned empty string")
	}
	
	// Default version should be "dev" if not set by build flags
	if version == "dev" {
		t.Log("Version is 'dev' (expected for test environment)")
	}
}

func TestGetBuildInfo(t *testing.T) {
	info := GetBuildInfo()
	
	// Test that all fields are populated
	if info.Version == "" {
		t.Error("BuildInfo.Version is empty")
	}
	
	if info.GitCommit == "" {
		t.Error("BuildInfo.GitCommit is empty")
	}
	
	if info.BuildDate == "" {
		t.Error("BuildInfo.BuildDate is empty")
	}
	
	if info.GoVersion == "" {
		t.Error("BuildInfo.GoVersion is empty")
	}
	
	if info.Platform == "" {
		t.Error("BuildInfo.Platform is empty")
	}
	
	// Test that GoVersion starts with "go"
	if !strings.HasPrefix(info.GoVersion, "go") {
		t.Errorf("GoVersion %q does not start with 'go'", info.GoVersion)
	}
	
	// Test that Platform contains expected format "OS/ARCH"
	if !strings.Contains(info.Platform, "/") {
		t.Errorf("Platform %q does not contain '/' separator", info.Platform)
	}
	
	// Verify that runtime values are correctly populated
	expectedGoVersion := runtime.Version()
	if info.GoVersion != expectedGoVersion {
		t.Errorf("GoVersion = %q, want %q", info.GoVersion, expectedGoVersion)
	}
	
	expectedPlatform := runtime.GOOS + "/" + runtime.GOARCH
	if info.Platform != expectedPlatform {
		t.Errorf("Platform = %q, want %q", info.Platform, expectedPlatform)
	}
}

func TestBuildInfoDefaults(t *testing.T) {
	info := GetBuildInfo()
	
	// Test default values when ldflags are not set
	if info.Version == "" || info.Version == "dev" {
		// This is expected in test environment
		t.Logf("Version is default value: %s", info.Version)
	}
	
	if info.GitCommit == "" || info.GitCommit == "unknown" {
		// This is expected in test environment
		t.Logf("GitCommit is default value: %s", info.GitCommit)
	}
	
	if info.BuildDate == "" || info.BuildDate == "unknown" {
		// This is expected in test environment  
		t.Logf("BuildDate is default value: %s", info.BuildDate)
	}
}

func TestInfoString(t *testing.T) {
	tests := []struct {
		name     string
		info     Info
		contains []string
		notContains []string
	}{
		{
			name: "complete build info",
			info: Info{
				Version:   "v1.2.3",
				GitCommit: "abc123def456",
				BuildDate: "2024-01-15T10:30:00Z",
				GoVersion: "go1.21.0",
				Platform:  "linux/amd64",
			},
			contains: []string{
				"pyinit v1.2.3",
				"Commit: abc123def456",
				"Built: 2024-01-15T10:30:00Z",
				"Go: go1.21.0",
				"Platform: linux/amd64",
			},
			notContains: []string{
				"unknown",
			},
		},
		{
			name: "minimal build info with unknowns",
			info: Info{
				Version:   "dev",
				GitCommit: "unknown",
				BuildDate: "unknown",
				GoVersion: "go1.21.0",
				Platform:  "darwin/arm64",
			},
			contains: []string{
				"pyinit dev",
				"Go: go1.21.0",
				"Platform: darwin/arm64",
			},
			notContains: []string{
				"Commit:",
				"Built:",
			},
		},
		{
			name: "partial build info",
			info: Info{
				Version:   "v0.1.0-beta",
				GitCommit: "deadbeef",
				BuildDate: "unknown",
				GoVersion: "go1.20.5",
				Platform:  "windows/amd64",
			},
			contains: []string{
				"pyinit v0.1.0-beta",
				"Commit: deadbeef",
				"Go: go1.20.5",
				"Platform: windows/amd64",
			},
			notContains: []string{
				"Built:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.info.String()
			
			// Check that expected strings are present
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("String() result does not contain %q\nResult:\n%s", expected, result)
				}
			}
			
			// Check that unexpected strings are not present
			for _, notExpected := range tt.notContains {
				if strings.Contains(result, notExpected) {
					t.Errorf("String() result contains %q but shouldn't\nResult:\n%s", notExpected, result)
				}
			}
		})
	}
}

func TestInfoStringFormat(t *testing.T) {
	info := Info{
		Version:   "v1.0.0",
		GitCommit: "abc123",
		BuildDate: "2024-01-01",
		GoVersion: "go1.21.0",
		Platform:  "linux/amd64",
	}
	
	result := info.String()
	
	// Check that the result starts with application name and version
	if !strings.HasPrefix(result, "pyinit v1.0.0") {
		t.Errorf("String() should start with 'pyinit v1.0.0', got: %s", result)
	}
	
	// Check that the result contains newlines (multiline format)
	if !strings.Contains(result, "\n") {
		t.Error("String() should be multiline format")
	}
	
	// Count lines
	lines := strings.Split(result, "\n")
	expectedLines := 5 // App+Version, Commit, Built, Go, Platform
	if len(lines) != expectedLines {
		t.Errorf("String() should have %d lines, got %d: %v", expectedLines, len(lines), lines)
	}
}

func TestInfoStringWithEmptyValues(t *testing.T) {
	info := Info{
		Version:   "",
		GitCommit: "",
		BuildDate: "",
		GoVersion: "",
		Platform:  "",
	}
	
	result := info.String()
	
	// Should still contain the app name even with empty version
	if !strings.Contains(result, "pyinit") {
		t.Errorf("String() should contain app name even with empty values, got: %s", result)
	}
}

func TestInfoStringConsistency(t *testing.T) {
	// Test that calling String() multiple times returns same result
	info := GetBuildInfo()
	
	first := info.String()
	second := info.String()
	
	if first != second {
		t.Error("String() method should return consistent results")
	}
}

// Test interaction between GetBuildInfo and String
func TestBuildInfoStringIntegration(t *testing.T) {
	info := GetBuildInfo()
	result := info.String()
	
	// Should contain the app name
	if !strings.Contains(result, Name) {
		t.Errorf("String() should contain app name %q", Name)
	}
	
	// Should contain version (even if it's "dev")
	if !strings.Contains(result, info.Version) {
		t.Errorf("String() should contain version %q", info.Version)
	}
	
	// Should contain Go version
	if !strings.Contains(result, info.GoVersion) {
		t.Errorf("String() should contain Go version %q", info.GoVersion)
	}
	
	// Should contain platform
	if !strings.Contains(result, info.Platform) {
		t.Errorf("String() should contain platform %q", info.Platform)
	}
}

// Test constants
func TestConstants(t *testing.T) {
	if Name == "" {
		t.Error("Name constant should not be empty")
	}
	
	if Name != "pyinit" {
		t.Errorf("Name constant = %q, want %q", Name, "pyinit")
	}
}

// Test that global variables have reasonable defaults
func TestGlobalVariables(t *testing.T) {
	// These variables might be set by ldflags, so we just check they exist
	if len(Version) == 0 {
		t.Error("Version variable should not be empty")
	}
	
	if len(GitCommit) == 0 {
		t.Error("GitCommit variable should not be empty")
	}
	
	if len(BuildDate) == 0 {
		t.Error("BuildDate variable should not be empty")
	}
	
	t.Logf("Version: %s", Version)
	t.Logf("GitCommit: %s", GitCommit)
	t.Logf("BuildDate: %s", BuildDate)
}

// Benchmark the version functions
func BenchmarkGetVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetVersion()
	}
}

func BenchmarkGetBuildInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetBuildInfo()
	}
}

func BenchmarkInfoString(b *testing.B) {
	info := GetBuildInfo()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = info.String()
	}
}

// Test cross-platform behavior
func TestPlatformInfo(t *testing.T) {
	info := GetBuildInfo()
	
	// Platform should match current runtime
	expectedPlatform := runtime.GOOS + "/" + runtime.GOARCH
	if info.Platform != expectedPlatform {
		t.Errorf("Platform = %q, want %q", info.Platform, expectedPlatform)
	}
	
	// Platform should contain valid OS values
	validOS := []string{"linux", "darwin", "windows", "freebsd", "openbsd", "netbsd"}
	osFound := false
	for _, os := range validOS {
		if strings.Contains(info.Platform, os) {
			osFound = true
			break
		}
	}
	if !osFound {
		t.Errorf("Platform %q does not contain a recognized OS", info.Platform)
	}
	
	// Platform should contain valid architecture values
	validArch := []string{"amd64", "386", "arm64", "arm", "ppc64", "ppc64le", "mips", "mipsle"}
	archFound := false
	for _, arch := range validArch {
		if strings.Contains(info.Platform, arch) {
			archFound = true
			break
		}
	}
	if !archFound {
		t.Errorf("Platform %q does not contain a recognized architecture", info.Platform)
	}
}