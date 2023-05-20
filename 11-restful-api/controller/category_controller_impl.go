package controller

import (
	"net/http"
	"strconv"

	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/helper"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/model/web"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/service"
	"github.com/julienschmidt/httprouter"
)

type CategoryControllerImpl struct {
	CategoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) CategoryController {
	return &CategoryControllerImpl{
		CategoryService: categoryService,
	}
}

func (controller *CategoryControllerImpl) Create(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	categoryCreateRequest := web.CategoryCreateRequest{}

	helper.ReadFromRequestBody(request, &categoryCreateRequest)

	categoryResponse := controller.CategoryService.Create(request.Context(), categoryCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   categoryResponse,
	}
	helper.SendToWriter(writer, webResponse)

}

func (controller *CategoryControllerImpl) Update(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	categoryUpdateRequest := web.CategoryUpdateRequest{}

	helper.ReadFromRequestBody(request, &categoryUpdateRequest)

	categoryId, err := strconv.Atoi(param.ByName("categoryId"))
	helper.PanifIfError(err)
	categoryUpdateRequest.Id = categoryId

	categoryResponse := controller.CategoryService.Update(request.Context(), categoryUpdateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponse,
	}

	helper.SendToWriter(writer, webResponse)
}

func (controller *CategoryControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {

	categoryId, err := strconv.Atoi(param.ByName("categoryId"))
	helper.PanifIfError(err)
	id := categoryId

	controller.CategoryService.Delete(request.Context(), id)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.SendToWriter(writer, webResponse)
}

func (controller *CategoryControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	categoryId, err := strconv.Atoi(param.ByName("categoryId"))
	helper.PanifIfError(err)
	id := categoryId

	categoryResponse := controller.CategoryService.FindById(request.Context(), id)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponse,
	}

	helper.SendToWriter(writer, webResponse)
}

func (controller *CategoryControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	categoryResponse := controller.CategoryService.FindAll(request.Context())
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Ok",
		Data:   categoryResponse,
	}

	helper.SendToWriter(writer, webResponse)
}
