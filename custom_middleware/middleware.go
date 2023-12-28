package custom_middleware

import (
	"context"
	"echo_mongo/db"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func ValidateJWT(c echo.Context) int {
	errorVariables := godotenv.Load()
	if errorVariables != nil {
		http.Error(c.Response(), "error parsing .env file", http.StatusUnauthorized)
		return 0
	}
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	header := c.Request().Header.Get("Authorization")
	if len(header) == 0 {
		return 0
	}
	splitBearer := strings.Split(header, " ")
	if len(splitBearer) != 2 {
		return 0
	}
	splitToken := strings.Split(splitBearer[1], ".")
	if len(splitToken) != 3 {
		return 0
	}
	tk := strings.TrimSpace(splitBearer[1])
	token, err := jwt.Parse(tk, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: ")

		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secretKey, nil
	})
	if err != nil {
		return 0
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := bson.M{}
		if err := db.UsersCollection.FindOne(context.TODO(), bson.M{
			"email": claims["email"],
		}).Decode(&user); err != nil {
			return 0
		}
		return 1
	} else {
		return 0
	}
}
