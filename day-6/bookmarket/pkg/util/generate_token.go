package util

import (
	"bookmarket/internal/app/users/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Claim struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(ctx echo.Context, user models.Users) (string, error) {
	// create expired time 1 day
	// Create 1 day
	expirationTime := time.Now().AddDate(0, 0, 1)

	// Create clain for payload token
	claim := Claim{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// Signed token with secret key
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	// If success, return token
	return signedToken, nil
}
