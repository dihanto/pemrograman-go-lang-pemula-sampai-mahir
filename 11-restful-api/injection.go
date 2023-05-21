package main

import (
	"net/http"

	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/app"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/controller"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/middleware"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/repository"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/service"
	"github.com/go-playground/validator"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
)

var categorySet = wire.NewSet(
	repository.NewCategoryRepository,
	wire.Bind(new(repository.CategoryRepository), new(*repository.CategoryRepositoryImpl)),
	service.NewCategoryService,
	wire.Bind(new(service.CategoryService), new(*service.CategoryServiceImpl)),
	controller.NewCategoryController,
	wire.Bind(new(controller.CategoryController), new(*controller.CategoryControllerImpl)),
)

func InitializedServer() *http.Server {
	wire.Build(
		app.NewDb,
		validator.New,
		categorySet,
		app.NewRouter,
		wire.Bind(new(http.Handler), new(*httprouter.Router)),
		middleware.NewAuthMiddleware,
		NewServer,
	)
	return nil

}
