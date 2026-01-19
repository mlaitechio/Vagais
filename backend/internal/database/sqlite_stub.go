//go:build !sqlite

package database

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
)

// InitializeSQLite sets up SQLite database connection (stub when SQLite is not available)
func InitializeSQLite(cfg *config.Config) (*gorm.DB, error) {
	return nil, fmt.Errorf("SQLite support not compiled in. Use -tags sqlite to enable SQLite support")
}



