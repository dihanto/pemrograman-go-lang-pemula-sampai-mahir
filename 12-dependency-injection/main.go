package main

import (
	"net/http"

	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/helper"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/middleware"
	_ "github.com/go-sql-driver/mysql"
)

func NewServer(authMiddleware *middleware.AuthMiddleware) *http.Server {
	return &http.Server{
		Addr:    "localhost:3000",
		Handler: authMiddleware,
	}
}

func main() {
	server := InitializedServer()
	err := server.ListenAndServe()
	helper.PanifIfError(err)
}
