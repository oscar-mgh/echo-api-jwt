package jwt

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerateJWT(id, email, name string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		panic("error loading .env file")
	}
	key := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"name":  name,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	encodedToken, err := token.SignedString(key)

	return encodedToken, err
}
