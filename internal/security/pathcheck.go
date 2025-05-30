// Package pathcheck provides path validation helpers for security checks.
package pathcheck

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// ValidatePath checks that the path is absolute and does not contain ".." segments.
func ValidatePath(path string) error {
	if !filepath.IsAbs(path) {
		return errors.New("path must be absolute")
	}

	if hasDotDotSegment(path) {
		return errors.New("path contains forbidden '..' segment")
	}

	return nil
}

func hasDotDotSegment(path string) bool {
	parts := strings.Split(path, string(os.PathSeparator))
	return slices.Contains(parts, "..")
}
