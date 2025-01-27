package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRefreshToken(t *testing.T) {
	token, err := MakeRefreshToken()
	assert.NoError(t, err, "MakeRefreshToken should not error")
	assert.NotEqual(t, token, "", "Token should not be empty")
}