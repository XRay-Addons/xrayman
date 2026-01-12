package launchctl

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/exec"
)

const launchctl = "/bin/launchctl"

// create service
func createService(domain, plistLocation string) error {
	_, _, err := exec.Run(launchctl, []string{
		"bootstrap", domain, plistLocation,
	})
	if err != nil {
		return fmt.Errorf("%w: bootstrap: %v", errdefs.ErrService, err)
	}
	return nil
}

// not exist is not an error
func removeService(domain, service string) error {
	_, stderr, err := exec.Run(launchctl, []string{
		"bootout", filepath.Join(domain, service),
	})

	const notExistsStderr = "No such process"
	if err == nil || strings.Contains(stderr, notExistsStderr) {
		return nil
	}
	return fmt.Errorf("%w: bootout: %v", errdefs.ErrService, err)
}

func startService(domain, service string) error {
	_, _, err := exec.Run(launchctl, []string{
		"kickstart", "-k", filepath.Join(domain, service),
	})
	if err != nil {
		return fmt.Errorf("%w: kickstart: %v", errdefs.ErrService, err)
	}
	return nil
}

// not exist or not running is not an error
func stopService(domain, service string) error {
	_, stderr, err := exec.Run(launchctl, []string{
		"kill", "TERM", filepath.Join(domain, service),
	})

	const alreadyStoppedStderr = "No process to signal."
	if err != nil && !strings.Contains(stderr, alreadyStoppedStderr) {
		return fmt.Errorf("%w: kill: %v", errdefs.ErrService, err)
	}
	return nil
}

// not exists is not an error, returns stopped
func getServiceStatus(domain, service string) (string, error) {
	stdout, _, err := exec.Run(launchctl, []string{
		"print", filepath.Join(domain, service),
	})

	if err != nil {
		return "", fmt.Errorf("%w: print: %v", errdefs.ErrService, err)
	}

	return extractStatusString(stdout), nil
}

const statusRegex = `(?m)^[ \t]*state = (.+?)$`

// parse state = something and extract 'something'
func extractStatusString(s string) string {
	re := regexp.MustCompile(statusRegex)
	match := re.FindStringSubmatch(s)
	if len(match) == 2 {
		return match[1]
	}
	return ""
}
