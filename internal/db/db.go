package db

import (
	"fmt"
	"myapp/internal/config"
	"myapp/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB struct encapsulates a GORM database connection.
// It provides methods for database operations and management.
type DB struct {
	*gorm.DB // Embeds a GORM database connection pointer.
}

// Connect establishes a connection to the PostgreSQL database using the provided configuration.
// It returns a *DB instance on success, or an error if the connection fails.
func Connect(cfg *config.Config) (*DB, error) {
	// Construct the database connection string (DSN) from the configuration.
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	// Open a new GORM database connection using the PostgreSQL driver and DSN.
	// Set the logger to display informational logs.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		// Return nil and the error if the connection fails.
		return nil, err
	}

	// Return a new DB instance wrapping the GORM database connection.
	return &DB{DB: db}, nil
}

// Close closes the underlying SQL database connection.
// It returns an error if closing the connection fails.
func (db *DB) Close() error {
	// Retrieve the underlying SQL database connection from the GORM connection.
	sqlDB, err := db.DB.DB()
	if err != nil {
		// Return the error if retrieving the connection fails.
		return err
	}

	// Close the SQL database connection.
	return sqlDB.Close()
}

// Migrate automatically migrates the database schema based on the defined models.
// It returns an error if the migration fails.
// Currently migrates the User model.
func Migrate(db *DB) error {
	// Use GORM's AutoMigrate to create or update the User table.
	return db.AutoMigrate(&model.User{})
}
