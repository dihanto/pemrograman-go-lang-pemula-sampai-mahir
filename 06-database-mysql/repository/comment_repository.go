package repository

import (
	"context"

	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/06-database-mysql/entity"
)

type CommentRepository interface {
	Insert(ctx context.Context, comment entity.Comment) (entity.Comment, error)
	FindById(ctx context.Context, id int) (entity.Comment, error)
	FindAll(ctx context.Context) ([]entity.Comment, error)
}
