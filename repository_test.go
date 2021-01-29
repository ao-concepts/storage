package storage_test

import (
	"testing"

	"github.com/ao-concepts/storage"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type entity struct {
	gorm.Model
	Value string
}

func mockEntity(value string) *entity {
	return &entity{
		Value: value,
	}
}

type entityRepository struct {
	storage.Repository
}

func (t *entityRepository) getAll(tx *storage.Transaction, entries *[]entity) error {
	return tx.Gorm().Find(entries).Error
}

func (t *entityRepository) getAllUnscoped(tx *storage.Transaction, entries *[]entity) error {
	return tx.Gorm().Unscoped().Find(entries).Error
}

var er entityRepository

func checkCommit(c *storage.Controller, assert *assert.Assertions) {
	var entries []entity

	assert.Nil(c.UseTransaction(func(tx *storage.Transaction) (err error) {
		return er.getAll(tx, &entries)
	}))
	assert.Len(entries, 1)

	entry := entries[0]
	assert.NotEmpty(entry.ID)
	assert.Equal(entry.Value, "test-value")
}

func checkRollback(c *storage.Controller, assert *assert.Assertions) {
	var entries []entity

	assert.Nil(c.UseTransaction(func(tx *storage.Transaction) (err error) {
		return er.getAll(tx, &entries)
	}))
	assert.Len(entries, 0)
}

func checkLen(len int, tx *storage.Transaction, assert *assert.Assertions) {
	var entries []entity
	assert.Nil(er.getAll(tx, &entries))
	assert.Len(entries, len)
}

func createMockStorageController() *storage.Controller {
	c, _ := storage.New(sqlite.Open(":memory:"), nil)
	c.Gorm().AutoMigrate(&entity{})
	return c
}

func TestRepository_Insert(t *testing.T) {
	assert := assert.New(t)
	c := createMockStorageController()

	assert.Nil(c.UseTransaction(func(tx *storage.Transaction) (err error) {
		return er.Insert(tx, mockEntity("test-value"))
	}))
	checkCommit(c, assert)
}

func TestRepository_Update(t *testing.T) {
	assert := assert.New(t)
	c := createMockStorageController()
	data := mockEntity("test-value")

	assert.Nil(c.UseTransaction(func(tx *storage.Transaction) (err error) {
		return er.Insert(tx, data)
	}))

	assert.Nil(c.UseTransaction(func(tx *storage.Transaction) (err error) {
		data.Value = "Updated string"
		return er.Update(tx, data)
	}))

	// check update
	var entries []entity

	assert.Nil(c.UseTransaction(func(tx *storage.Transaction) (err error) {
		return er.getAll(tx, &entries)
	}))
	assert.Len(entries, 1)

	entry := entries[0]
	assert.NotEmpty(entry.ID)
	assert.Equal("Updated string", entry.Value)
}

func TestRepository_Delete(t *testing.T) {
	assert := assert.New(t)
	c := createMockStorageController()
	data := mockEntity("test-value")

	assert.Nil(c.UseTransaction(func(tx *storage.Transaction) (err error) {
		return er.Insert(tx, data)
	}))
	checkCommit(c, assert)

	assert.Nil(c.UseTransaction(func(tx *storage.Transaction) (err error) {
		return er.Delete(tx, data)
	}))

	var entries []entity

	assert.Nil(c.UseTransaction(func(tx *storage.Transaction) (err error) {
		return er.getAllUnscoped(tx, &entries)
	}))
	assert.Len(entries, 1)
	checkRollback(c, assert)
}

func TestRepository_Remove(t *testing.T) {
	assert := assert.New(t)
	c := createMockStorageController()
	data := mockEntity("test-value")

	assert.Nil(c.UseTransaction(func(tx *storage.Transaction) (err error) {
		return er.Insert(tx, data)
	}))
	checkCommit(c, assert)

	assert.Nil(c.UseTransaction(func(tx *storage.Transaction) (err error) {
		return er.Remove(tx, data)
	}))

	var entries []entity

	assert.Nil(c.UseTransaction(func(tx *storage.Transaction) (err error) {
		return er.getAllUnscoped(tx, &entries)
	}))
	assert.Len(entries, 0)
}
