package handlers

import (
	"context"
	"echo_mongo/custom_middleware"
	"echo_mongo/db"
	"echo_mongo/dto"
	"echo_mongo/jwt"
	"echo_mongo/validations"
	"encoding/json"
	"net/http"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Protected(c echo.Context) error {
	if custom_middleware.ValidateJWT(c) == 0 {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "unauthorized",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusUnauthorized)
		return json.NewEncoder(c.Response()).Encode(sr)
	}
	sr := dto.StandardResponse{
		Status: "ok",
		Msg:    "access granted to protected route, jwt validated!",
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(sr)
}

func Login(c echo.Context) error {
	loginDto := dto.LoginDto{}
	if err := json.NewDecoder(c.Request().Body).Decode(&loginDto); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "unexpected error ocurred",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if len(loginDto.Email) == 0 || validations.ValidEmail.FindStringSubmatch(loginDto.Email) == nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "invalid email",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if !validations.ValidPassword(loginDto.Password) {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "the password must be between 8 and 22 characters long, have an upper case letter and a number",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	user := bson.M{}
	if err := db.UsersCollection.FindOne(context.TODO(), bson.M{
		"email": loginDto.Email,
	}).Decode(&user); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "email is not registered",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	bytePassword := []byte(loginDto.Password)
	dbPassword := []byte(user["password"].(string))
	err := bcrypt.CompareHashAndPassword(dbPassword, bytePassword)

	if err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "invalid credentials",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	} else {
		strObjId := user["_id"].(primitive.ObjectID).Hex()
		jwtKey, err := jwt.GenerateJWT(strObjId, user["email"].(string), user["name"].(string))

		if err != nil {
			sr := dto.StandardResponse{
				Status: "error",
				Msg:    "something went wrong generating jwt token",
			}
			c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(sr)
		} else {
			sr := dto.LoginResponseDto{
				Name:  user["name"].(string),
				Token: jwtKey,
			}
			c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
			c.Response().WriteHeader(http.StatusOK)
			return json.NewEncoder(c.Response()).Encode(sr)
		}
	}
}

func Register(c echo.Context) error {
	userDto := dto.UserDto{}
	if err := json.NewDecoder(c.Request().Body).Decode(&userDto); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "unexpected error ocurred",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if len(userDto.Name) == 0 {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "name is required",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if len(userDto.Email) == 0 || validations.ValidEmail.FindStringSubmatch(userDto.Email) == nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "invalid email",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if len(userDto.Phone) == 0 {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "invalid phone",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if !validations.ValidPassword(userDto.Password) {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "the password must be between 8 and 22 characters long, have an upper case letter and a number",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	user := bson.M{}
	if err := db.UsersCollection.FindOne(context.TODO(), bson.M{"email": userDto.Email}).Decode(&user); err == nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "email address is already registered",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	doc := bson.D{
		{Key: "name", Value: userDto.Name},
		{Key: "slug", Value: slug.Make(userDto.Name)},
		{Key: "phone", Value: userDto.Phone},
		{Key: "email", Value: userDto.Email},
		{Key: "password", Value: string(hash)},
	}
	db.UsersCollection.InsertOne(context.TODO(), doc)

	sr := dto.StandardResponse{
		Status: "created",
		Msg:    "user created successfully",
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(sr)
}
