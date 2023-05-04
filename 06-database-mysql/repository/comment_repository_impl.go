package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/06-database-mysql/entity"
)

type commentRepositoryImpl struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepositoryImpl{DB: db}
}

func (cri *commentRepositoryImpl) Insert(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	script := "INSERT INTO comments(email, comment) VALUES(?,?)"
	result, err := cri.DB.ExecContext(ctx, script, comment.Email, comment.Comment)
	if err != nil {
		panic(err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	comment.Id = int32(lastId)
	return comment, err
}

func (cri *commentRepositoryImpl) FindById(ctx context.Context, id int) (entity.Comment, error) {
	script := "SELECT id, email, comment FROM comments WHERE id=? LIMIT 1"
	rows, err := cri.DB.QueryContext(ctx, script, id)
	comment := entity.Comment{}
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&comment.Id, &comment.Email, comment.Comment)
		return comment, nil
	} else {
		return comment, errors.New("Id" + strconv.Itoa(int(id)) + " Not Found")
	}
}

func (cri *commentRepositoryImpl) FindAll(ctx context.Context) ([]entity.Comment, error) {
	script := "Select id, email, comment FROM comments"
	rows, err := cri.DB.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var comments []entity.Comment
	for rows.Next() {
		comment := entity.Comment{}
		rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		comments = append(comments, comment)
	}
	return comments, nil
}
