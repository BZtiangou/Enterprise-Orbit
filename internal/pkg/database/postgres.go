package database

import (
	"fmt"
	"log"
	"os"
	"time"

	config "enterprise-orbit/configs"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgreSQL(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	if os.Getenv("USE_SQLITE") == "true" || cfg.Host == "" {
		return NewSQLite()
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println("PostgreSQL connection failed, falling back to SQLite")
		return NewSQLite()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to PostgreSQL successfully")
	return db, nil
}

func NewSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("enterprise_orbit.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
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

	log.Println("Connected to SQLite successfully")
	return db, nil
}
