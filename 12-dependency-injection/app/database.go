package app

import (
	"database/sql"
	"time"

	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/helper"
)

func NewDb() *sql.DB {

	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/restfulapi")
	helper.PanifIfError(err)

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)
	return db
}
