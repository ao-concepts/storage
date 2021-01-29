package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_Commit(t *testing.T) {
	assert := assert.New(t)
	c := createMockStorageController()
	tx, _ := c.Begin()

	assert.Nil(er.Insert(tx, mockEntity("test-value")))
	checkLen(1, tx, assert)
	assert.Nil(tx.Commit())
	checkCommit(c, assert)
}

func TestTransaction_Rollback(t *testing.T) {
	assert := assert.New(t)
	c := createMockStorageController()
	tx, _ := c.Begin()

	assert.Nil(er.Insert(tx, mockEntity("value-value")))
	checkLen(1, tx, assert)
	assert.Nil(tx.Rollback())
	checkRollback(c, assert)
}

func TestTransaction_Gorm(t *testing.T) {
	assert := assert.New(t)
	c := createMockStorageController()
	tx, _ := c.Begin()

	assert.NotNil(tx.Gorm())
}
