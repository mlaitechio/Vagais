package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/models"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

// Initialize sets up the database connection (PostgreSQL or SQLite)
func Initialize(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if cfg.DatabaseType == "sqlite" {
		// Use SQLite for local development
		db, err = InitializeSQLite(cfg)
		if err != nil {
			return nil, err
		}
	} else {
		// Use PostgreSQL for production
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PostgreSQL database: %v", err)
		}
	}

	// Set connection pool settings (only for PostgreSQL)
	if cfg.DatabaseType != "sqlite" {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("failed to get database instance: %v", err)
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	// Auto-migrate models
	if err := migrateModels(db); err != nil {
		return nil, fmt.Errorf("failed to migrate models: %v", err)
	}

	DB = db
	return db, nil
}

// InitializeRedis sets up Redis connection with graceful fallback
func InitializeRedis(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed, continuing without Redis: %v", err)
		return nil, nil // Return nil instead of error for graceful fallback
	}

	RedisClient = client
	return client, nil
}

// migrateModels performs database migrations
func migrateModels(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Organization{},
		&models.User{},
		&models.Agent{},
		&models.Review{},
		&models.Execution{},
		&models.LLMProvider{},
		&models.PasswordResetToken{},
	)
}

// IsRedisAvailable checks if Redis is available
func IsRedisAvailable() bool {
	return RedisClient != nil
}

// CloseConnections closes all database connections
func CloseConnections() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}
