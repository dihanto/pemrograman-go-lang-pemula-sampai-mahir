package httprouter

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestRouterPatternNammedParammeter(t *testing.T) {

	router := httprouter.New()
	router.GET("/products/:id/items/:itemid", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		itemid := p.ByName("itemid")
		text := "Product " + id + " items " + itemid
		fmt.Fprint(w, text)
	})

	request := httptest.NewRequest("GET", "http://localhost:2323/products/1/items/1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, "Product 1 items 1", string(body))
}
func TestRouterPatternCatchAllParammeter(t *testing.T) {

	router := httprouter.New()
	router.GET("/images/*image", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		image := p.ByName("image")
		text := "images : " + image
		fmt.Fprint(w, text)
	})

	request := httptest.NewRequest("GET", "http://localhost:2323/images/small/picture.png", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, "images : /small/picture.png", string(body))
}
