package test

import (
	"testing"

	"github.com/dihanto/pemrograman-go-lang-pemula-sampai-mahir/11-restful-api/simple"
	"github.com/stretchr/testify/assert"
)

func TestSimpleServiceError(t *testing.T) {
	simpleService, err := simple.InitializedService(true)
	assert.NotNil(t, err)
	assert.Nil(t, simpleService)
}

func TestSimpleServiceSuccess(t *testing.T) {
	simpleService, err := simple.InitializedService(false)
	assert.Nil(t, err)
	assert.NotNil(t, simpleService)
}
