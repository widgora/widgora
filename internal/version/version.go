// Package version holds build metadata injected at link time with -ldflags -X.
package version

var (
	// Version is the release version, e.g. v0.1.0.
	Version = "dev"
	// Commit is the short git commit the binary was built from.
	Commit = "none"
	// Date is the commit time of that build in RFC 3339.
	Date = "unknown"
)
