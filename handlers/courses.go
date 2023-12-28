package handlers

import (
	"context"
	"echo_mongo/db"
	"echo_mongo/dto"
	"encoding/json"
	"net/http"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindAllCourses(c echo.Context) error {
	pipeline := []bson.M{
		{"$match": bson.M{}},
		{"$lookup": bson.M{"from": "categories", "localField": "category_id", "foreignField": "_id", "as": "category"}},
		{"$sort": bson.M{"_id": -1}},
	}

	courses, err := db.CourseCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		panic(err)
	}
	results := []bson.M{}
	if err = courses.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(results)
}

func FindCourseById(c echo.Context) error {
	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	result := bson.M{}
	if err := db.CourseCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "course not found",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusNotFound)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(result)
}

func CreateCourse(c echo.Context) error {
	courseDto := dto.CourseDto{}
	if err := json.NewDecoder(c.Request().Body).Decode(&courseDto); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "unexpected error ocurred",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if len(courseDto.Name) == 0 {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "name is required",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if len(courseDto.Description) < 5 {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "invalid description",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if courseDto.Price == 0 {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "invalid price",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	categoryID, _ := primitive.ObjectIDFromHex(courseDto.CategoryID)
	doc := bson.D{
		{Key: "name", Value: courseDto.Name},
		{Key: "price", Value: courseDto.Price},
		{Key: "description", Value: courseDto.Description},
		{Key: "category_id", Value: categoryID},
		{Key: "slug", Value: slug.Make(courseDto.Name)},
	}

	db.CourseCollection.InsertOne(context.TODO(), doc)
	sr := dto.StandardResponse{
		Status: "created",
		Msg:    "course created successfully",
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(sr)
}

func UpdateCourse(c echo.Context) error {
	courseDto := dto.CourseDto{}
	if err := json.NewDecoder(c.Request().Body).Decode(&courseDto); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "unexpected error ocurred",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if len(courseDto.Name) == 0 {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "name is required",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if len(courseDto.Description) < 5 {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "invalid description",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if courseDto.Price == 0 {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "invalid price",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	result := bson.M{}
	if err := db.CourseCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "unexpected error ocurred",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	doc := make(map[string]any)
	doc["name"] = courseDto.Name
	doc["price"] = courseDto.Price
	doc["description"] = courseDto.Description
	doc["slug"] = slug.Make(courseDto.Name)
	updateString := bson.M{"$set": doc}
	db.CourseCollection.UpdateOne(context.TODO(), bson.M{"_id": bson.M{"$eq": objID}}, updateString)

	sr := dto.StandardResponse{
		Status: "updated",
		Msg:    "course updated successfully",
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(sr)
}

func DeleteCourse(c echo.Context) error {
	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	result := bson.M{}
	if err := db.CourseCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "unexpected error ocurred",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	db.CourseCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	sr := dto.StandardResponse{
		Status: "deleted",
		Msg:    "category deleted successfully",
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusNoContent)
	return json.NewEncoder(c.Response()).Encode(sr)
}
