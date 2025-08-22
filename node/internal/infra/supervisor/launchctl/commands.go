package launchctl

import (
	"context"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/XRay-Addons/xrayman/node/internal/infra/exec"
)

const launchctl = "/bin/launchctl"

// create service
func createService(ctx context.Context, domain, plistLocation string) error {
	_, _, err := exec.Run(ctx, launchctl, []string{
		"bootstrap", domain, plistLocation,
	})
	if err != nil {
		return err
	}
	return nil
}

// not exist is not an error
func removeService(ctx context.Context, domain, service string) error {
	_, stderr, err := exec.Run(ctx, launchctl, []string{
		"bootout", filepath.Join(domain, service),
	})

	const notExistsStderr = "No such process"
	if err == nil || strings.Contains(stderr, notExistsStderr) {
		return nil
	}
	return err
}

func startService(ctx context.Context, domain, service string) error {
	_, _, err := exec.Run(ctx, launchctl, []string{
		"kickstart", "-k", filepath.Join(domain, service),
	})
	if err != nil {
		return err
	}
	return nil
}

// not exist or not running is not an error
func stopService(ctx context.Context, domain, service string) error {
	_, stderr, err := exec.Run(ctx, launchctl, []string{
		"kill", "TERM", filepath.Join(domain, service),
	})

	const alreadyStoppedStderr = "No process to signal."
	if err != nil && !strings.Contains(stderr, alreadyStoppedStderr) {
		return err
	}
	return nil
}

// not exists is not an error, returns stopped
func getServiceStatus(ctx context.Context, domain, service string) (string, error) {
	stdout, _, err := exec.Run(ctx, launchctl, []string{
		"print", filepath.Join(domain, service),
	})

	if err != nil {
		return "", err
	}

	return extractStatusString(stdout), nil
}

const statusRegex = `(?m)^[ \t]*state = (.+?)$`

// parse state = something and extract 'something'
func extractStatusString(s string) string {
	re := regexp.MustCompile(statusRegex)
	match := re.FindStringSubmatch(s)
	const statusParts = 2
	if len(match) == statusParts {
		return match[1]
	}
	return ""
}
