package storage

import (
	"gorm.io/gorm"
)

// Transaction for a storage action
type Transaction struct {
	gormTx *gorm.DB
}

// Commit a transaction
func (tx *Transaction) Commit() (err error) {
	return tx.gormTx.Commit().Error
}

// Rollback a transaction
func (tx *Transaction) Rollback() (err error) {
	return tx.gormTx.Rollback().Error
}

// Gorm returns the gorm transaction
func (tx *Transaction) Gorm() *gorm.DB {
	return tx.gormTx
}
