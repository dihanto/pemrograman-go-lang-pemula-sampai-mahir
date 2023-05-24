package service

import (
	"context"
	"database/sql"

	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/exception"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/helper"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/model/domain"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/model/web"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/repository"
	"github.com/go-playground/validator"
)

type CategoryServiceImpl struct {
	CategoryRepo repository.CategoryRepository
	DB           *sql.DB
	Validate     *validator.Validate
}

func NewCategoryService(categoryRepo repository.CategoryRepository, db *sql.DB, validate *validator.Validate) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		CategoryRepo: categoryRepo,
		DB:           db,
		Validate:     validate,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse {
	err := service.Validate.Struct(request)
	helper.PanifIfError(err)

	tx, err := service.DB.Begin()
	helper.PanifIfError(err)
	defer helper.RollbackOrCommit(tx)

	category := domain.Category{
		Name: request.Name,
	}

	category = service.CategoryRepo.Save(ctx, tx, category)
	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse {
	err := service.Validate.Struct(request)
	helper.PanifIfError(err)

	tx, err := service.DB.Begin()
	helper.PanifIfError(err)
	defer helper.RollbackOrCommit(tx)

	category, err := service.CategoryRepo.FindById(ctx, tx, request.Id)
	if err != nil {
		panic(exception.NewNotfoundError(err.Error()))
	}

	category.Name = request.Name

	category = service.CategoryRepo.Update(ctx, tx, category)

	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) {
	tx, err := service.DB.Begin()
	helper.PanifIfError(err)
	defer helper.RollbackOrCommit(tx)

	category, err := service.CategoryRepo.FindById(ctx, tx, categoryId)
	if err != nil {
		panic(exception.NewNotfoundError(err.Error()))
	}

	service.CategoryRepo.Delete(ctx, tx, category)
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) web.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanifIfError(err)
	defer helper.RollbackOrCommit(tx)

	category, err := service.CategoryRepo.FindById(ctx, tx, categoryId)
	if err != nil {
		panic(exception.NewNotfoundError(err.Error()))
	}

	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []web.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanifIfError(err)
	defer helper.RollbackOrCommit(tx)
	categories := service.CategoryRepo.FindAll(ctx, tx)

	return helper.ToCategoryResponses(categories)
}
