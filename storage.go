package storage

import (
	"fmt"

	"gorm.io/gorm"
)

// Controller of storage
type Controller struct {
	db *gorm.DB
}

type Config struct {
	MaxOpenConnections     int
	MaxIdleConnections     int
	SkipDefaultTransaction bool
}

// New storage controller consructor.
// Pass nil to `log` to use the default logger of gorm.
// default settings:
//    * SetMaxIdleConns: 0 = unlimited
// 	  * SetMaxOpenConns: 0 = unlimited
func New(dialector gorm.Dialector, config *Config, log Logger) (c *Controller, err error) {
	if dialector == nil {
		return nil, fmt.Errorf("Cannot use a storage controller without a dialector")
	}

	if config == nil {
		config = &Config{
			MaxOpenConnections:     0,
			MaxIdleConnections:     0,
			SkipDefaultTransaction: true,
		}
	}

	cfg := &gorm.Config{}
	if log != nil {
		cfg.Logger = log.CreateGormLogger()
	}

	db, err := gorm.Open(dialector, cfg)
	if err != nil {
		return nil, err
	}

	db.SkipDefaultTransaction = config.SkipDefaultTransaction

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(config.MaxOpenConnections)

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
