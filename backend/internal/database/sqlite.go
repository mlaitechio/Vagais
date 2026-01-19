//go:build sqlite

package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mlaitechio/vagais/internal/config"
)

// InitializeSQLite sets up SQLite database connection
func InitializeSQLite(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("agais.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %v", err)
	}
	return db, nil
}



