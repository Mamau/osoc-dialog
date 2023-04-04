package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func CreateAccessToken(id int, name string, secret string) (accessToken string, err error) {
	expires := time.Now().Add(time.Hour)
	claims := &CustomClaims{
		Name: name,
		ID:   id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func CreateRefreshToken(id int, secret string) (refreshToken string, err error) {
	expires := time.Now().Add(time.Hour)
	claimsRefresh := &RefreshClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func ExtractIDFromToken(requestToken string, secret string) (int, error) {
	token, err := jwt.ParseWithClaims(requestToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*CustomClaims)

	if !ok && !token.Valid {
		return 0, fmt.Errorf("invalid Token")
	}

	return claims.ID, nil
}
