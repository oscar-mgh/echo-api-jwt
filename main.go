package main

import (
	"echo_mongo/db"
	"echo_mongo/handlers"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Static("/public", "public")

	db.TestConnection()

	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Println(errEnv)
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Go - Echo API Working!!!")
	})

	categories := e.Group("/api/v1/categories")
	categories.GET("", handlers.FindAllCategories)
	categories.GET("/:id", handlers.FindCategoryById)
	categories.POST("", handlers.CreateCategory)
	categories.PUT("/:id", handlers.UpdateCategory)
	categories.DELETE("/:id", handlers.DeleteCategory)

	courses := e.Group("/api/v1/courses")
	courses.GET("", handlers.FindAllCourses)
	courses.GET("/thumbnails/:id", handlers.FindAllThumbnails)
	courses.DELETE("/thumbnails/:id", handlers.DeleteThumbnails)
	courses.GET("/:id", handlers.FindCourseById)
	courses.PUT("/:id", handlers.UpdateCourse)
	courses.POST("", handlers.CreateCourse)
	courses.POST("/:id", handlers.UploadThumbnail)
	courses.DELETE("/:id", handlers.DeleteCourse)

	security := e.Group("/api/v1/users")
	security.POST("/register", handlers.Register)
	security.POST("/login", handlers.Login)
	security.GET("/protected", handlers.Protected)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://**", "https://**"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
