package postgres

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
	DBurl    string
}

func SetupDatabase(db *gorm.DB, models ...interface{}) error {
	err := db.AutoMigrate(models...)
	return err
}

func New(config *Config) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	if os.Getenv("ENV") == "staging" || os.Getenv("ENV") == "production" {
		db, err = gorm.Open(postgres.Open(config.DBurl), &gorm.Config{})
	} else {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
