package repository

import "github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/03-unit-test/entity"

type CategoryRepository interface {
	FindById(id string) *entity.Category
}
