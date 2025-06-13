package furniture

import "mis/storage/models"

type FurnitureRepo interface {
	GetFurnitureByID(id uint) (models.Furniture, error)
	GetFurniture() ([]models.Furniture, error)
	AddFurniture(furniture models.Furniture) (uint, error)
	UpdateFurniture(furniture models.Furniture) error
}
