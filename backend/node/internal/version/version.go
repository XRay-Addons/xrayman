package version

import "fmt"

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func String() string {
	return fmt.Sprintf(
		`XRay Node.
	Version: %s
	Commit: %s
	BuildTime: %s`,
		Version, Commit, BuildTime)
}
