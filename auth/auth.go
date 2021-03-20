package auth

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

func CheckAuth(authToken string) (bool, string) {
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return false, ""
	}

	user := token.Claims.(jwt.MapClaims)["name"].(string)

	if user == "" {
		return false, ""
	}

	return true, user
}
