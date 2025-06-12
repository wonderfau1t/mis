package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"mis/storage/models"
)

type Storage struct {
	db *gorm.DB
}

func SetupStorage() (*Storage, error) {
	const fn = "storage.postgresql.SetupStorage"

	dsn := "host=db user=postgres password=postgres dbname=mis port=5432"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	db.AutoMigrate(
		&models.Status{}, &models.Category{}, &models.Furniture{},
		&models.Order{}, &models.PartOfOrder{})

	return &Storage{db: db}, nil
}
