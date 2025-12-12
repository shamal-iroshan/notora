package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// CreateAccess generates a signed JWT access token containing the user's ID
// and an expiration time. The token uses HS256 (HMAC SHA-256) algorithm.
//
// Parameters:
//   - secret: signing key
//   - userID: ID of the authenticated user
//   - ttlSeconds: token lifetime in seconds
//
// Used for:
//   - Short-lived access tokens in HTTP-only cookies
func CreateAccess(secret []byte, userID int64, ttlSeconds int) (string, error) {
	claims := jwtlib.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(ttlSeconds) * time.Second).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// Parse validates and parses a JWT access token.
// It returns the claims if the token is valid.
//
// Security measures included:
//   - Ensures token was signed with HS256
//   - Validates signature using the provided secret
//
// Parameters:
//   - secret: signing key
//   - tokenString: raw JWT string from cookie
//
// Returns:
//   - claims (MapClaims)
//   - error if invalid or expired
func Parse(secret []byte, tokenString string) (jwtlib.MapClaims, error) {
	claims := jwtlib.MapClaims{}

	parsedToken, err := jwtlib.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwtlib.Token) (interface{}, error) {

			// ------------------------------------------------------------------
			// SECURITY CHECK:
			// Ensure token uses expected algorithm.
			// Prevents "alg=none" and algorithm substitution attacks.
			// ------------------------------------------------------------------
			if t.Method != jwtlib.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
