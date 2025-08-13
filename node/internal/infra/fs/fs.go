package fs

import (
	"fmt"
	"os"
)

func AccessFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("fs: access file: %w", err)
	}
	return !info.IsDir(), nil
}

func AccessDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("fs: access dir: %w", err)
	}
	return info.IsDir(), nil
}
