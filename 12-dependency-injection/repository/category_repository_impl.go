package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/helper"
	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/model/domain"
)

type CategoryRepositoryImpl struct {
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
}

func (cri *CategoryRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category {
	sql := "insert into category(name) values (?)"
	result, err := tx.ExecContext(ctx, sql, category.Name)
	helper.PanifIfError(err)

	id, err := result.LastInsertId()
	helper.PanifIfError(err)

	category.Id = int(id)
	return category
}

func (cri *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category {
	sql := "update category set name = ? where id = ?"
	_, err := tx.ExecContext(ctx, sql, category.Name, category.Id)
	helper.PanifIfError(err)

	return category
}

func (cri *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, category domain.Category) {
	sql := "delete from category where id = ?"
	_, err := tx.ExecContext(ctx, sql, category.Id)
	helper.PanifIfError(err)
}

func (cri *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, categoryId int) (domain.Category, error) {
	sql := "select id, name from category where id = ?"
	rows, err := tx.QueryContext(ctx, sql, categoryId)
	helper.PanifIfError(err)
	defer rows.Close()

	category := domain.Category{}
	if rows.Next() {
		err = rows.Scan(&category.Id, &category.Name)
		helper.PanifIfError(err)
		return category, nil
	} else {
		return category, errors.New("category not found")
	}
}

func (cri *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Category {
	sql := "select id, name from category"
	rows, err := tx.QueryContext(ctx, sql)
	helper.PanifIfError(err)
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		category := domain.Category{}
		err = rows.Scan(&category.Id, &category.Name)
		helper.PanifIfError(err)
		categories = append(categories, category)
	}
	return categories
}
