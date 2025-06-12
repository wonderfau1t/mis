package models

import (
	"time"
)

type Status struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique" json:"name"`
}

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique" json:"name"`

	Furniture []Furniture `gorm:"foreignKey:CategoryID" json:"-"`
}

type Furniture struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique"`
	CategoryID  uint
	Price       uint
	Description *string
	Photo       *string

	Category Category
}

type Order struct {
	ID                uint `gorm:"primaryKey"`
	CreatedAt         time.Time
	CustomerFirstName string
	CustomerLastName  string
	CustomerPhone     string
	Sum               uint
	StatusID          uint

	PartsOfOrder []PartOfOrder `gorm:"foreignKey:OrderID"`
	Status       Status
}

type PartOfOrder struct {
	ID          uint `gorm:"primaryKey"`
	OrderID     uint
	FurnitureID uint
	StatusID    uint

	Order     Order
	Furniture Furniture
	Status    Status
}
