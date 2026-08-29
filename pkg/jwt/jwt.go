package jwt

import (
	"errors"
	"strings"

	jwt2 "github.com/golang-jwt/jwt/v5"
)

func CreateToken(secret []byte, claims *UserClaims) (string, error) {
	return jwt2.NewWithClaims(jwt2.SigningMethodHS512, claims).SignedString(secret)
}

func ParseToken(tokenString string, secret []byte) (*UserClaims, error) {
	if strings.Count(tokenString, ".") != 2 {
		return nil, errors.New("token contains an invalid number of segments")
	}
	token, err := jwt2.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt2.Token) (interface{}, error) {
		return secret, nil
	})
	if token == nil {
		return nil, errors.New("token contains an invalid token")
	}
	var claims *UserClaims
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("token contains an invalid claims")
	}
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}
