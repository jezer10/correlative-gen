// internal/database/db.go
package database

import (
	"github.com/emidiaz3/event-driven-server/config"
	"gorm.io/gorm"
)

// GetDB returns the global DB connection
func GetDB() *gorm.DB {
	return config.DB
}
