package networkconfig_test

import (
	"testing"

	"github.com/gologames/go-mvp/internal/networkconfig"
	"github.com/stretchr/testify/assert"
)

func TestValidateIdentifierCorrect(t *testing.T) {
	t.Parallel()

	identifiers := []string{
		"",
		"name",
		"a",
		"dev1",
		"host.local",
		"o000000000",
		"-",
		"-0-0-0-0-0-",
		"a.0.-",
		".",
		"..",
		"用户",
		"____test____",
	}

	for _, identifier := range identifiers {
		err := networkconfig.ValidateIdentifier(identifier)
		assert.NoError(t, err)
	}
}

func TestValidateIdentifierIncorrect(t *testing.T) {
	t.Parallel()

	identifiers := []string{
		"1",
		"5guys",
		"123",
		"user!",
		"@mail",
		"a b",
		"a/b",
		"a\tb",
		"a\nb",
	}

	for _, identifier := range identifiers {
		err := networkconfig.ValidateIdentifier(identifier)
		assert.Error(t, err)
	}
}
