package util

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken generates a JWT RSA token with the user ID as part of the claims
func GenerateRSAToken(userID uint) (string, error) {
	prvKey, errPrvKey := os.ReadFile("cert/id_rsa")
	if errPrvKey != nil {
		log.Fatalln(errPrvKey)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = userID                        // Our custom data.
	claims["exp"] = now.Add(time.Hour * 1).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()                    // The time at which the token was issued.
	claims["nbf"] = now.Unix()                    // The time before which the token must be disregarded.

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

// VerifyToken verifies a token JWT RSA validate
func VerifyRSAToken(tokenString string) (jwt.MapClaims, error) {
	pubKey, errPubKey := os.ReadFile("cert/id_rsa.pub")
	if errPubKey != nil {
		log.Fatalln(errPubKey)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims, nil
}
