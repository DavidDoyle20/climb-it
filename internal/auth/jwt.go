package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	expiresIn := time.Duration(time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
		&jwt.RegisteredClaims{
			Issuer: "climbit", 
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()), 
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject: userID.String(),
		})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// parse the token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func (token *jwt.Token) (interface{}, error) {
		// validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { 
			return nil, jwt.ErrSignatureInvalid
		}
		// return the secret key for validation
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	// type assert the claims
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}

	// validate the issuer
	if claims.Issuer != "climbit" {
		return uuid.Nil, jwt.ErrTokenInvalidIssuer
	}

	// validate the time based claims
	if time.Now().After(claims.ExpiresAt.Time) {
		return uuid.Nil, jwt.ErrTokenExpired
	}
	if time.Now().Before(claims.IssuedAt.Time) {
		return uuid.Nil, jwt.ErrTokenNotValidYet
	}

	// convert the subject back into uuid
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, jwt.ErrTokenInvalidSubject
	}

	return userID, nil
	}
