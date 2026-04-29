package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chirpy-access",
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	mySigningKey := []byte(tokenSecret)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", fmt.Errorf("error signing key")
	}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return uuid.Nil, fmt.Errorf("validate error: %v", err)
	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		uuid, err := uuid.Parse(claims.Subject)
		if err != nil {
			log.Fatal("unknown uuid")
		}
		return uuid, nil
	} else {
		log.Fatal("unknown claims type, cannot proceed")
	}
	return uuid.Nil, fmt.Errorf("Unknown error")
}
