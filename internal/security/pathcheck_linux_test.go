//go:build linux

package pathcheck_test

import (
	"testing"

	pathcheck "github.com/gologames/go-mvp/internal/security"
	"github.com/stretchr/testify/assert"
)

func TestValidatePath_Correct(t *testing.T) {
	t.Parallel()

	validPaths := []string{
		"/etc/passwd",
		"/usr/local/bin",
		"/var/log/nginx/access.log",
		"/home/user/dir/file.txt",
	}

	for _, path := range validPaths {
		assert.NoError(t, pathcheck.ValidatePath(path))
	}
}

func TestValidatePath_NotAbsolute(t *testing.T) {
	t.Parallel()

	invalidPaths := []string{
		"relative/path",
		"./file",
		"..",
		"folder/file.txt",
		"../etc/passwd",
		"../tmp",
	}

	for _, path := range invalidPaths {
		assert.Error(t, pathcheck.ValidatePath(path), path)
	}
}

func TestValidatePath_HasDotDotSegment(t *testing.T) {
	t.Parallel()

	withDotDot := []string{
		"/usr/local/../bin",
		"/etc/../var/log",
		"/var/../tmp/../lib",
		"/tmp/..",
	}

	for _, path := range withDotDot {
		assert.Error(t, pathcheck.ValidatePath(path), path)
	}
}
