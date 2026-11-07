package jwt

import (
	"errors"
	"fmt"
	"strings"

	jwt2 "github.com/golang-jwt/jwt/v5"
)

func CreateToken(secret []byte, claims *UserClaims) (string, error) {
	return jwt2.NewWithClaims(jwt2.SigningMethodHS512, claims).SignedString(secret)
}
func ParseToken(tokenString string, secret []byte) (*UserClaims, error) {

	tokenString = strings.TrimSpace(tokenString)

	if strings.Count(tokenString, ".") != 2 {
		return nil, errors.New(
			"token contains an invalid number of segments",
		)
	}

	token, err := jwt2.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt2.Token) (interface{}, error) {

			if token.Method != jwt2.SigningMethodHS512 {
				return nil, fmt.Errorf(
					"unexpected signing method: %s",
					token.Method.Alg(),
				)
			}

			return secret, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf(
			"parse token: %w",
			err,
		)
	}

	if token == nil || !token.Valid {
		return nil, errors.New(
			"token is invalid",
		)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New(
			"token contains invalid claims",
		)
	}

	return claims, nil
}
