package utils

import (
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
)

type TokenClaims struct {
	// Meta data
	Claim interface{} `json:"claim"`

	// Inherit from standard claims
	jwt.StandardClaims
}

// GenerateJWTToken
func GenerateJWTToken(privateKeydata []byte, claim TokenClaims, expireOffsetHour int64) (string, error) {

	privateKey, keyErr := jwt.ParseECPrivateKeyFromPEM(privateKeydata)
	if keyErr != nil {
		log.Fatalf("unable to parse private key: %s", keyErr.Error())
	}

	method := jwt.GetSigningMethod(jwt.SigningMethodES256.Name)

	session, err := jwt.NewWithClaims(method, claim).SignedString(privateKey)
	return session, err
}

// ValidateToken
func ValidateToken(keydata []byte, token string) (jwt.MapClaims, error) {

	publicKey, keyErr := jwt.ParseECPublicKeyFromPEM(keydata)
	if keyErr != nil {
		return nil, keyErr
	}

	parsed, parseErr := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if parseErr != nil {
		return nil, parseErr
	}

	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		log.Printf("Claims JWT: %v", claims)
		return claims, nil
	} else {
		return nil, fmt.Errorf("Token claim is not valid!")
	}
}
