package handlers

import (
	"context"
	"echo_mongo/db"
	"echo_mongo/dto"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UploadThumbnail(c echo.Context) error {
	file, err := c.FormFile("photo")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	extension := strings.Split(file.Filename, ".")[1]
	time := strings.Split(time.Now().String(), " ")
	photo := string(time[4][6:14]) + "." + extension
	archive := "public/uploads/courses_thumbnails/" + photo

	dst, err := os.Create(archive)
	if err != nil {
		return err
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	result := bson.M{}
	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	doc := bson.D{
		{Key: "name", Value: photo},
		{Key: "course_id", Value: objID},
	}

	if err := db.CourseCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "course not found",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusNotFound)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	db.CourseThumbnailCollection.InsertOne(context.TODO(), doc)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	sr := dto.StandardResponse{
		Status: "ok",
		Msg:    "thumbnail uploaded successfully",
	}
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(sr)
}

func FindAllThumbnails(c echo.Context) error {
	result := bson.M{}
	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := db.CourseCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "course not found",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusNotFound)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	thumbnails, err := db.CourseThumbnailCollection.Find(context.TODO(), bson.M{"course_id": objID})
	if err != nil {
		panic(err)
	}
	results := []bson.M{}
	if err := thumbnails.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(results)
}

func DeleteThumbnails(c echo.Context) error {
	result := bson.M{}
	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := db.CourseThumbnailCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "thumbnails not found",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusNotFound)
		return json.NewEncoder(c.Response()).Encode(sr)
	}
	path := "public/uploads/courses_thumbnails/" + result["name"].(string)
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
	db.CourseThumbnailCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	sr := dto.StandardResponse{
		Status: "no content",
		Msg:    "thumbnails deleted successfully",
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(sr)
}
