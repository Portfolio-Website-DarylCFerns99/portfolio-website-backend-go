package security

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
)

func getSecretKey() []byte {
	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		key = "YOUR_DEFAULT_SECRET_KEY_CHANGE_THIS"
	}
	return []byte(key)
}

func getExpireMinutes() int {
	minStr := os.Getenv("ACCESS_TOKEN_EXPIRE_MINUTES")
	if minStr == "" {
		return 30
	}
	min, err := strconv.Atoi(minStr)
	if err != nil {
		return 30
	}
	return min
}

// CreateAccessToken generates a new JWT token signed with the secret key.
func CreateAccessToken(data map[string]interface{}) (string, error) {
	expireMinutes := getExpireMinutes()
	expirationTime := time.Now().UTC().Add(time.Duration(expireMinutes) * time.Minute)

	claims := jwt.MapClaims{
		"exp": expirationTime.Unix(),
	}

	for k, v := range data {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}

// DecodeToken verifies a JWT token and returns the claims payload.
func DecodeToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure algorithm matches HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return getSecretKey(), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, ErrInvalidToken
}
