package storage_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ao-concepts/storage"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/logger"
)

type gormLogger struct {
}

func (l *gormLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
}

func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
}

func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
}

type testLogger struct {
}

func (l *testLogger) CreateGormLogger() logger.Interface {
	return &gormLogger{}
}

func TestNew(t *testing.T) {
	assert := assert.New(t)

	// missing dialector
	_, err := storage.New(nil, nil)
	assert.NotNil(err)

	// missing logger
	_, err = storage.New(sqlite.Open(":memory:"), nil)
	assert.Nil(err)

	// custom logger
	_, err = storage.New(sqlite.Open(":memory:"), &testLogger{})
	assert.Nil(err)

	// simple
	c, err := storage.New(sqlite.Open(":memory:"), nil)
	assert.Nil(err)
	assert.NotNil(c)

	// wrong connections
	_, err = storage.New(mysql.Open("not-a-connection-string"), nil)
	assert.NotNil(err)
}

func TestController_Gorm(t *testing.T) {
	assert := assert.New(t)

	c, err := storage.New(sqlite.Open(":memory:"), nil)
	assert.Nil(err)
	assert.NotNil(c.Gorm())
}

func TestController_Begin(t *testing.T) {
	assert := assert.New(t)
	c := createMockStorageController()

	tx, err := c.Begin()
	assert.Nil(err)
	assert.NotNil(tx)
}

func TestController_UseTransaction(t *testing.T) {
	assert := assert.New(t)
	c := createMockStorageController()

	// Rollback
	err := c.UseTransaction(func(tx *storage.Transaction) (err error) {
		assert.Nil(er.Insert(tx, mockEntity("test-value")))
		checkLen(1, tx, assert)
		return fmt.Errorf("Some test error")
	})
	assert.Equal(err, fmt.Errorf("Some test error"))
	checkRollback(c, assert)

	// Commit
	err = c.UseTransaction(func(tx *storage.Transaction) (err error) {
		assert.Nil(er.Insert(tx, mockEntity("test-value")))
		checkLen(1, tx, assert)

		return nil
	})
	assert.Nil(err)
	checkCommit(c, assert)
}
