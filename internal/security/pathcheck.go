// Package pathcheck provides path validation helpers for security checks.
package pathcheck

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// ValidatePath checks that the path is absolute and does not contain ".." segments.
func ValidatePath(path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path must be absolute: %q", path)
	}

	if hasDotDotSegment(path) {
		return fmt.Errorf("path contains forbidden '..' segment: %q", path)
	}

	return nil
}

func hasDotDotSegment(path string) bool {
	parts := strings.Split(path, string(os.PathSeparator))
	return slices.Contains(parts, "..")
}
