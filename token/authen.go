package token

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// สร้าง JWT Token โดยการผสม Secret
func GenerateToken(username string, secret string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{username},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	})

	signedToken, err := t.SignedString([]byte(secret))
	if err != nil {
		log.Println("error signing key")
		return signedToken, err
	}

	return signedToken, nil
}

func VerifyToken(tokenString string, secret string) (*jwt.Token, error) {
	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Return the verified token
	return t, nil
}
