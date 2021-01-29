package storage

import "gorm.io/gorm/logger"

// Logger interface
type Logger interface {
	CreateGormLogger() logger.Interface
}
