package database_mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO customers(id, name) VALUES('ujang32', 'UUjang')"
	_, err := db.ExecContext(ctx, script)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new customer")
}
func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name FROM customers"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		fmt.Println("id", id)
		fmt.Println("name", name)

	}
}

func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name, email, balance, rating, birth_data, married, create_at FROM customers"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating float64
		var createdAt time.Time
		var birthDate sql.NullTime
		var married bool
		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}
		fmt.Println("==============")
		fmt.Println("id", id)
		fmt.Println("name", name)
		if email.Valid {
			fmt.Println("email", email.String)
		}
		fmt.Println("balance", balance)
		fmt.Println("rating", rating)
		if birthDate.Valid {
			fmt.Println("birth date", birthDate.Time)

		}
		fmt.Println("married", married)
		fmt.Println("created at", createdAt)

	}
}

func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin'; #"
	password := "salah"

	script := "SELECT username  FROM user WHERE username ='" + username + "' AND password ='" + password + "' LIMIT 1"
	fmt.Println(script)
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Suksen Login", username)
	} else {
		fmt.Println("Gagal Login")
	}
}

func TestSqlInjectionSafe(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin"
	password := "admin"

	script := "SELECT username  FROM user WHERE username =? AND password =? LIMIT 1"
	fmt.Println(script)
	rows, err := db.QueryContext(ctx, script, username, password)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Suksen Login", username)
	} else {
		fmt.Println("Gagal Login")
	}
}

func TestExecSqlParameter(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "dihanto23"
	password := "dihanto"

	script := "INSERT INTO user(username, password) VALUES(?, ?)"
	_, err := db.ExecContext(ctx, script, username, password)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new user")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "dihanto233@gmail.com"
	comment := "test comment"

	script := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	result, err := db.ExecContext(ctx, script, email, comment)
	if err != nil {
		panic(err)
	}
	lastid, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new comment with id", lastid)
}

func TestPrepareStetement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO comments(email, comment) VALUES(?, ?)"

	statement, err := db.PrepareContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	for i := 0; i < 10; i++ {
		email := "dihanto" + strconv.Itoa(i) + "@gmail.com"
		comment := "test comment"

		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}
		lastid, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("success comment with id ", lastid)
	}
}

func TestTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	script := "INSERT INTO comments(email, comment) VALUES(?, ?)"

	for i := 0; i < 10; i++ {
		email := "dihanto" + strconv.Itoa(i) + "@gmail.com"
		comment := "test comment"

		result, err := tx.ExecContext(ctx, script, email, comment)
		if err != nil {
			panic(err)
		}
		lastid, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("success comment with id ", lastid)
	}

	err = tx.Rollback()
	if err != nil {
		panic(err)
	}
}
