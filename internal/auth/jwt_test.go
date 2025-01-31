package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestMakeJWT(t *testing.T) {
	// test data
	userID := "testid"
	tokenSecret := "test_secret"
	expiresIn := time.Hour

	tokenString, err := MakeJWT(userID, tokenSecret)

	assert.NoError(t, err, "Should not return an error when creating a token")
	assert.NotEmpty(t, tokenString, "Should return a non-empty token string")

	// Parse the token to verify its contents
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(tokenString, &jwt.RegisteredClaims{})
	assert.NoError(t, err, "Should be able to parse the token")

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	assert.True(t, ok, "Should be able to assert claims type")
	assert.Equal(t, "climbit", claims.Issuer, "Issuer should be climbit")
	assert.Equal(t, userID, claims.Subject, "Subject should match userId")
	assert.WithinDuration(t, time.Now().UTC(), claims.IssuedAt.Time, time.Second, "IssuedAt should be close to now")
	assert.WithinDuration(t, time.Now().UTC().Add(expiresIn), claims.ExpiresAt.Time, time.Second, "ExpiresAt should be one hour from now")
}

func TestValidateJWT(t *testing.T) {
	// Test data
	userID := "testid"
	tokenSecret := "test_secret"
	expiresIn := time.Hour

	// Create a valid token
	tokenString, err := MakeJWT(userID, tokenSecret)
	assert.NoError(t, err, "Should not error when creating a token")

	// Test valid token validation
	validatedID, err := ValidateJWT(tokenString, tokenSecret)
	assert.NoError(t, err, "Should not error when validating a token")
	assert.Equal(t, userID, validatedID, "Should return the correct user id")

	// Test token with wrong secret
	_, err = ValidateJWT(tokenString, "wrong_secret")
	assert.Error(t, err, "Should error with incorrect secret")
	assert.ErrorIs(t, err, jwt.ErrSignatureInvalid, "Error should be due to signature invalid")

	// Test token with invalid claims (e.g. wrong issuer)
	wrongIssuerToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    "wrong_isser",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID,
	})
	wrongIssuerTokenString, err := wrongIssuerToken.SignedString([]byte(tokenSecret))
	assert.NoError(t, err, "Should not error signing token with wrong issuer")
	_, err = ValidateJWT(wrongIssuerTokenString, tokenSecret)
	assert.Error(t, err, "Should error with wrong issuer")
	assert.ErrorIs(t, err, jwt.ErrTokenInvalidIssuer, "Error should be due to invalid issuer")

	// Test malformed token
	_, err = ValidateJWT("malformed_token", tokenSecret)
	assert.Error(t, err, "Should error on malformed token")

}
