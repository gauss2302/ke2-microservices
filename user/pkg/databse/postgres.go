package database

import (
	"fmt"
	"log"
	"time"

	"github.com/gauss2302/testcommm/user/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg config.DataBaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 0; i < cfg.MaxRetries; i++ {
		db, err = gorm.Open(postgres.Open(cfg.URL), &gorm.Config{ // убрали :=
			Logger: logger.Default.LogMode(logger.Error),
		})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database, attempt %d/%d: %v", i+1, cfg.MaxRetries, err)
		time.Sleep(cfg.RetryInterval)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", cfg.MaxRetries, err)
	}
	log.Printf("Successfully connected to database: %s", cfg.URL)
	return db, nil
}
