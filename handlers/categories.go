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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAllCategories(c echo.Context) error {
	// categories, err := db.CategoryCollection.Find(context.TODO(), bson.D{})
	sortedCategories, err := db.CategoryCollection.Find(
		context.TODO(),
		bson.D{},
		options.Find().SetSort(bson.D{{Key: "name", Value: +1}}),
	)
	if err != nil {
		panic(err)
	}
	results := []bson.M{}
	if err = sortedCategories.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(results)
}

func FindCategoryById(c echo.Context) error {
	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	result := bson.M{}
	if err := db.CategoryCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "category not found",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusNotFound)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(result)
}

func CreateCategory(c echo.Context) error {
	categoryDto := dto.CategoryDto{}
	if err := json.NewDecoder(c.Request().Body).Decode(&categoryDto); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "unexpected error ocurred",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if len(categoryDto.Name) == 0 {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "name is required",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	doc := bson.D{{Key: "name", Value: categoryDto.Name}, {Key: "slug", Value: slug.Make(categoryDto.Name)}}
	db.CategoryCollection.InsertOne(context.TODO(), doc)
	sr := dto.StandardResponse{
		Status: "created",
		Msg:    "category created successfully",
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(sr)
}

func UpdateCategory(c echo.Context) error {
	categoryDto := dto.CategoryDto{}
	if err := json.NewDecoder(c.Request().Body).Decode(&categoryDto); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "unexpected error ocurred",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	if len(categoryDto.Name) == 0 {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "name is required",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	result := bson.M{}
	if err := db.CategoryCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "unexpected error ocurred",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	doc := make(map[string]any)
	doc["name"] = categoryDto.Name
	doc["slug"] = slug.Make(categoryDto.Name)
	updateString := bson.M{"$set": doc}
	db.CategoryCollection.UpdateOne(context.TODO(), bson.M{"_id": bson.M{"$eq": objID}}, updateString)

	sr := dto.StandardResponse{
		Status: "updated",
		Msg:    "category updated successfully",
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(sr)
}

func DeleteCategory(c echo.Context) error {
	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	result := bson.M{}
	if err := db.CategoryCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result); err != nil {
		sr := dto.StandardResponse{
			Status: "error",
			Msg:    "unexpected error ocurred",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(sr)
	}

	db.CategoryCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	sr := dto.StandardResponse{
		Status: "deleted",
		Msg:    "category deleted successfully",
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusNoContent)
	return json.NewEncoder(c.Response()).Encode(sr)
}
