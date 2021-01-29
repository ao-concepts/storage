package storage

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Controller of storage
type Controller struct {
	db *gorm.DB
}

// New storage controller consructor.
// Pass nil to `log` to use the default logger of gorm.
// default settings:
//    * SetMaxIdleConns: 10
// 	  * SetMaxOpenConns: 100
// 	  * SetConnMaxLifetime: 1h
func New(dialector gorm.Dialector, log Logger) (c *Controller, err error) {
	if dialector == nil {
		return nil, fmt.Errorf("Cannot use a storage controller without a dialector")
	}

	cfg := &gorm.Config{}
	if log != nil {
		cfg.Logger = log.CreateGormLogger()
	}

	db, err := gorm.Open(dialector, cfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Controller{
		db: db,
	}, nil
}

// Gorm returns the used gorm handler.
func (c *Controller) Gorm() *gorm.DB {
	return c.db
}

// Begin a transaction
func (c *Controller) Begin() (tx *Transaction, err error) {
	gormTx := c.db.Begin()

	return &Transaction{
		gormTx: gormTx,
	}, gormTx.Error
}

// HandlerFunc is to be handled by a storage transaction
type HandlerFunc func(tx *Transaction) (err error)

// UseTransaction executes a HandlerFunc within a storage transaction
func (c *Controller) UseTransaction(fn HandlerFunc) (err error) {
	tx, err := c.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	return tx.Commit()
}
