package main

import (
	"net/http"

	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/app"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/controller"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/exception"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/helper"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/middleware"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/repository"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/service"
	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

func main() {
	db := app.NewDb()
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)

	router := httprouter.New()
	router.POST("/api/categories", categoryController.Create)
	router.GET("/api/categories", categoryController.FindAll)
	router.GET("/api/categories/:categoryId", categoryController.FindById)
	router.PUT("/api/categories/:categoryId", categoryController.Update)
	router.DELETE("/api/categories/:categoryId", categoryController.Delete)

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}
	err := server.ListenAndServe()
	helper.PanifIfError(err)
}
