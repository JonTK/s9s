package version

import (
	"runtime"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	info := Get()

	if info.Version == "" {
		t.Error("Version should not be empty")
	}

	if info.GoVersion != runtime.Version() {
		t.Errorf("Expected Go version %s, got %s", runtime.Version(), info.GoVersion)
	}

	expectedPlatform := runtime.GOOS + "/" + runtime.GOARCH
	if info.Platform != expectedPlatform {
		t.Errorf("Expected platform %s, got %s", expectedPlatform, info.Platform)
	}
}

func TestString(t *testing.T) {
	// Test with commit hash
	Version = "1.0.0"
	Commit = "abc123def456"
	Date = "2026-01-21"

	info := Get()
	str := info.String()

	if !strings.Contains(str, "1.0.0") {
		t.Errorf("String should contain version, got: %s", str)
	}

	if !strings.Contains(str, "abc123d") {
		t.Errorf("String should contain short commit hash, got: %s", str)
	}

	// Test with unknown commit
	Commit = "unknown"
	info = Get()
	str = info.String()

	if !strings.Contains(str, "1.0.0") {
		t.Errorf("String should contain version even without commit, got: %s", str)
	}
}

func TestShort(t *testing.T) {
	Version = "1.2.3"
	info := Get()

	if info.Short() != "1.2.3" {
		t.Errorf("Expected short version 1.2.3, got %s", info.Short())
	}
}

func TestFull(t *testing.T) {
	Version = "1.0.0"
	Commit = "abc123"
	Date = "2026-01-21"
	BuiltBy = "goreleaser"

	info := Get()
	full := info.Full()

	expectedStrings := []string{
		"s9s version",
		"1.0.0",
		"abc123",
		"2026-01-21",
		"goreleaser",
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(full, expected) {
			t.Errorf("Full() should contain '%s', got:\n%s", expected, full)
		}
	}
}
