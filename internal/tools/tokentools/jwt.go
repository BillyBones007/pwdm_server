package tokentools

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

// JWTTools - the tools for working with jwt tokens.
type JWTTools struct {
	secretKey []byte
}

// NewJWTTools - constructor JWTTools.
func NewJWTTools() *JWTTools {
	key, err := keyGenerate()
	if err != nil {
		log.Fatal(err)
	}
	return &JWTTools{secretKey: key}
}

// CreateToken - create a new token, signed secret key. Accepts the parametr expAt -
// time duration in unix format.
// For example:
// expAt := time.Now().Add(time.Hour * 1).Unix()
// CreateToken(expAt) returns token with an expiration time one hour.
func (j *JWTTools) CreateToken(expAt int64, uuid string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expAt
	claims["uuid"] = uuid

	tokenStr, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// ParseUUID - parse uuid from token string.
func (j *JWTTools) ParseUUID(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if expTime.Before(time.Now()) {
			return "", fmt.Errorf("token is expired")
		}

		uuid := claims["uuid"]
		if uuid.(string) == "" {
			return "", fmt.Errorf("uuid field is empty")
		}

		return uuid.(string), nil
	} else {
		return "", fmt.Errorf("invalid token claims")
	}
}

// keyGenerate - generating the secret key. Used for user token generating.
func keyGenerate() ([]byte, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
