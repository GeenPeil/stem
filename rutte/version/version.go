package version

import "log"

// this variable will be set by the linker at compile-time
var (
	version string
)

func init() {
	if len(version) == 0 {
		log.Fatal("empty version string")
	}
}

// String returns the complete version string.
func String() string {
	return version
}
