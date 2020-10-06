package version

import (
	"fmt"

	"github.com/coreos/go-semver/semver"
)

const (
	versionMajor int64 = 0
	versionMinor int64 = 0
	versionPatch int64 = 1
)

var Version = semver.Version{
	Major: versionMajor,
	Minor: versionMinor,
	Patch: versionPatch,
}

func VersionString(name string) string {
	return fmt.Sprintf("%s version %d.%d.%d", name, versionMajor, versionMinor, versionPatch)
}
