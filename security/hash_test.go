package security_test

import (
	"testing"

	"github.com/izaakdale/lib/security"
	"github.com/stretchr/testify/assert"
)

func TestHashAndVerifyPassword(t *testing.T) {
	password := "TestPassword"

	hashed, err := security.HashPassword(password, 1)
	assert.NoError(t, err)

	match := security.VerifyPassword(password, hashed)
	assert.True(t, match)
	match2 := security.VerifyPassword("WrongPassword", hashed)
	assert.False(t, match2)
}
