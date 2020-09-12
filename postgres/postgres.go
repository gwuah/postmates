package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

type Config struct {
	Host    string
	User    string
	DBName  string
	SSLMode string
}

func SetupDatabase(db *gorm.DB, models ...interface{}) error {
	err := db.AutoMigrate(models...)
	return err
}

func New(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s", config.Host, config.User, config.DBName, config.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
