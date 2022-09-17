package lib

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func ParseTokenJWT(tokenHeader string) (string, error) {
	// If there is, create new variable with empty string value
	encodedToken := ""
	// Split token on header with white space
	arrayToken := strings.Split(tokenHeader, " ")
	// If length arrayToken is same the 2
	if len(arrayToken) == 2 {
		// Get arrayToken with index 1 / only token jwt
		encodedToken = arrayToken[1]
	}

	// Parse token
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	// Get payload token
	claim, ok := token.Claims.(jwt.MapClaims)
	// If not `ok` and token invalid
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	// Get payload `id`
	userId := fmt.Sprint(claim["user_id"])

	return userId, nil

}
