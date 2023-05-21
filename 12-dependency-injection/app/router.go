package app

import (
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/controller"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/exception"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(categoryController controller.CategoryController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/categories", categoryController.FindAll)
	router.GET("/api/categories/:categoryId", categoryController.FindById)
	router.POST("/api/categories", categoryController.Create)
	router.PUT("/api/categories/:categoryId", categoryController.Update)
	router.DELETE("/api/categories/:categoryId", categoryController.Delete)

	router.PanicHandler = exception.ErrorHandler

	return router
}
