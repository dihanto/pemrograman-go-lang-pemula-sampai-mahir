package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache/v2"
	"github.com/stretchr/testify/assert"
)

var engine = mustache.New("./template", ".mustache")
var app = fiber.New(fiber.Config{
	Views: engine,
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		c.Status(fiber.StatusInternalServerError)
		return c.SendString("Error: " + err.Error())
	},
})

func TestRoutingHelloWorld(t *testing.T) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	request := httptest.NewRequest("GET", "/", nil)
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)
	assert.Equal(t, "Hello, World!", string(bytes))
}

func TestCtx(t *testing.T) {
	app.Get("/hello", func(c *fiber.Ctx) error {
		name := c.Query("name", "Guest")

		return c.SendString("Hello, " + name)
	})

	request := httptest.NewRequest("GET", "/hello", nil)
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, Guest", string(bytes))

	request = httptest.NewRequest("GET", "/hello?name=Dihanto", nil)
	response, err = app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err = io.ReadAll(response.Body)

	assert.Nil(t, err)
	assert.Equal(t, "Hello, Dihanto", string(bytes))

}

func TestHttpRequest(t *testing.T) {
	app.Get("/request", func(c *fiber.Ctx) error {
		first := c.Get("first", "Guest")
		last := c.Cookies("last", "Guest")
		return c.SendString("Hello, " + first + " " + last)
	})

	request := httptest.NewRequest("GET", "/request", nil)
	request.Header.Set("first", "Hans")
	request.AddCookie(&http.Cookie{Name: "last", Value: "Dihanto"})
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	assert.Equal(t, "Hello, Hans Dihanto", string(bytes))

}

func TestRouteParameters(t *testing.T) {
	app.Get("/users/:userId/orders/:orderId", func(c *fiber.Ctx) error {
		userId := c.Params("userId")
		orderId := c.Params("orderId")
		return c.SendString("Get Order " + orderId + " of User " + userId)
	})

	request := httptest.NewRequest("GET", "/users/dihanto/orders/2", nil)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	assert.Equal(t, "Get Order 2 of User dihanto", string(bytes))
}

func TestFormRequest(t *testing.T) {
	app.Post("/hello", func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		return c.SendString("Hello, " + name)
	})

	body := strings.NewReader("name=Dihanto")
	request := httptest.NewRequest("POST", "/hello", body)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := app.Test(request)

	assert.Nil(t, err)

	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	assert.Equal(t, "Hello, Dihanto", string(bytes))
}

//go:embed source/contoh.txt
var contohFile []byte

func TestFormUpload(t *testing.T) {
	app.Post("/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		err = c.SaveFile(file, "./target/"+file.Filename)
		if err != nil {
			return err
		}
		return c.SendString("Upload Success")
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, _ := writer.CreateFormFile("file", "contoh.txt")
	file.Write(contohFile)
	writer.Close()

	request := httptest.NewRequest("POST", "/upload", body)

	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := app.Test(request)

	assert.Nil(t, err)

	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	assert.Equal(t, "Upload Success", string(bytes))
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestRequestBody(t *testing.T) {
	app.Post("/login", func(c *fiber.Ctx) error {
		var request LoginRequest
		body := c.Body()
		err := json.Unmarshal(body, &request)
		if err != nil {
			return err
		}
		return c.SendString("Hello, " + request.Username)
	})

	body := strings.NewReader(`{"username":"dihanto", "password":"secret"}`)

	request := httptest.NewRequest("POST", "/login", body)

	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)

	assert.Nil(t, err)

	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	assert.Equal(t, "Hello, dihanto", string(bytes))
}

type RegisterRequest struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
	Name     string `json:"name" xml:"name" form:"name"`
}

func TestBodyParser(t *testing.T) {
	app.Post("/register", func(c *fiber.Ctx) error {
		var request RegisterRequest
		err := c.BodyParser(&request)
		if err != nil {
			return err
		}
		return c.SendString("Register Success, " + request.Username)
	})
}

func TestBodyParserJSON(t *testing.T) {
	TestBodyParser(t)

	body := strings.NewReader(`{"username":"dihanto", "password":"secret", "name":"Hans Dihanto"}`)
	request := httptest.NewRequest("POST", "/register", body)
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Register Success, dihanto", string(bytes))
}

func TestBodyParserForm(t *testing.T) {
	TestBodyParser(t)

	body := strings.NewReader("username=dihanto&password=secret&name=Hans Dihanto")

	request := httptest.NewRequest("POST", "/register", body)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := app.Test(request)

	assert.Nil(t, err)

	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	assert.Equal(t, "Register Success, dihanto", string(bytes))
}

func TestBodyParserXML(t *testing.T) {
	TestBodyParser(t)

	body := strings.NewReader(`
	<user>
		<username>dihanto</username>
		<password>secret</password>
		<name>Hans Dihanto</name>
	</user>`)

	request := httptest.NewRequest("POST", "/register", body)

	request.Header.Set("Content-Type", "application/xml")

	response, err := app.Test(request)

	assert.Nil(t, err)

	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	assert.Equal(t, "Register Success, dihanto", string(bytes))
}

func TestResponseJDON(t *testing.T) {
	app.Get("/user", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	request := httptest.NewRequest("GET", "/user", nil)

	response, err := app.Test(request)

	assert.Nil(t, err)

	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	assert.Equal(t, "{\"message\":\"Hello, World!\"}", string(bytes))
}

func TestDownloadFile(t *testing.T) {
	app.Get("/download", func(c *fiber.Ctx) error {
		return c.Download("./source/contoh.txt", "contoh.txt")
	})

	request := httptest.NewRequest("GET", "/download", nil)

	response, err := app.Test(request)

	assert.Nil(t, err)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "attachment; filename=\"contoh.txt\"", response.Header.Get("Content-Disposition"))

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	assert.Equal(t, "this is sample file for upload", string(bytes))
}

func TestRoutingGroup(t *testing.T) {
	helloWorld := func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	}
	api := app.Group("/api")
	api.Get("/hello", helloWorld)
	api.Get("world", helloWorld)

	web := api.Group("/web")
	web.Get("/hello", helloWorld)
	web.Get("/world", helloWorld)

	request := httptest.NewRequest("GET", "/api/hello", nil)

	response, err := app.Test(request)

	assert.Nil(t, err)

	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	assert.Equal(t, "Hello, World!", string(bytes))
}

func TestStaticFile(t *testing.T) {
	app.Static("/public", "./source")

	request := httptest.NewRequest("GET", "/public/contoh.txt", nil)
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	assert.Equal(t, "this is sample file for upload", string(bytes))
}

func TestErrorHandling(t *testing.T) {
	app.Get("/error", func(c *fiber.Ctx) error {
		return errors.New("Something went wrong")
	})

	request := httptest.NewRequest("GET", "/error", nil)

	response, err := app.Test(request)

	assert.Nil(t, err)

	assert.Equal(t, 500, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	assert.Equal(t, "Error: Something went wrong", string(bytes))
}

func TestView(t *testing.T) {
	app.Get("/view", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"title":   "Hello, World!",
			"header":  "Hello, Header!",
			"content": "Hello, Content!",
		})
	})

	request := httptest.NewRequest("GET", "/view", nil)

	response, err := app.Test(request)

	assert.Nil(t, err)

	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)

	assert.Nil(t, err)

	assert.Contains(t, string(bytes), "Hello, World!")
	assert.Contains(t, string(bytes), "Hello, Header!")
	assert.Contains(t, string(bytes), "Hello, Content!")

}

func TestClient(t *testing.T) {
	client := fiber.AcquireClient()
	defer fiber.ReleaseClient(client)

	agent := client.Get("https://example.com")
	status, response, errors := agent.String()
	assert.Nil(t, errors)
	assert.Equal(t, 200, status)
	assert.Contains(t, response, "Example Domain")
}
