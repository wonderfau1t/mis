package furniture

import "mis/storage/models"

type FurnitureRepo interface {
	GetFurniture() ([]models.Furniture, error)
	AddFurniture(furniture models.Furniture) (uint, error)
	UpdateFurniture(furniture models.Furniture) error
}
