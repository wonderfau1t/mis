package storage

import (
	"errors"
	"gorm.io/gorm"
	"mis/storage/models"
)

func (s *Storage) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := s.db.Find(&categories).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return categories, nil
		}
		return nil, err
	}
	return categories, nil
}

func (s *Storage) GetStatuses() ([]models.Status, error) {
	var statuses []models.Status
	if err := s.db.Find(&statuses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return statuses, nil
		}
		return nil, err
	}
	return statuses, nil
}

func (s *Storage) GetFurniture() ([]models.Furniture, error) {
	var furniture []models.Furniture
	if err := s.db.Preload("Category").Find(&furniture).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return furniture, nil
		}
		return nil, err
	}
	return furniture, nil
}

func (s *Storage) AddFurniture(furniture models.Furniture) (uint, error) {
	result := s.db.Create(&furniture)
	if result.Error != nil {
		return 0, result.Error
	}
	return furniture.ID, nil
}

func (s *Storage) UpdateFurniture(furniture models.Furniture) error {
	result := s.db.Save(&furniture)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *Storage) AddOrder(order models.Order) (uint, error) {
	result := s.db.Create(&order)
	if result.Error != nil {
		return 0, result.Error
	}
	return order.ID, nil
}

func (s *Storage) GetFurnitureByID(id uint) (models.Furniture, error) {
	var furniture models.Furniture
	if err := s.db.Preload("Category").First(&furniture, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return furniture, gorm.ErrRecordNotFound
		}
		return furniture, err
	}
	return furniture, nil
}

func (s *Storage) UpdateOrder(order models.Order) error {
	err := s.db.Where("order_id = ?", order.ID).Delete(&models.PartOfOrder{}).Error
	if err != nil {
		return err
	}
	result := s.db.Save(&order)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *Storage) GetOrders() ([]models.Order, error) {
	var orders []models.Order
	if err := s.db.Debug().Preload("PartsOfOrder").Preload("PartsOfOrder.Furniture").Find(&orders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return orders, nil
		}
		return nil, err
	}
	return orders, nil
}

func (s *Storage) GetOrderByID(orderID uint) (models.Order, error) {
	var order models.Order
	if err := s.db.Preload("PartsOfOrder").Preload("PartsOfOrder.Furniture").First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return order, gorm.ErrRecordNotFound
		}
		return order, err
	}
	return order, nil
}

func (s *Storage) UpdateOrderStatus(orderID uint, statusID uint) error {
	var order models.Order
	if err := s.db.Preload("PartsOfOrder").First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	for _, part := range order.PartsOfOrder {
		if part.StatusID != 3 {
			return errors.New("cannot change order status while parts of order are not completed")
		}
	}

	order.StatusID = statusID
	result := s.db.Save(&order)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) UpdatePartOfOrderStatus(partID uint, statusID uint) error {
	var part models.PartOfOrder
	if err := s.db.Preload("Order").First(&part, partID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	if statusID == 2 {
		if part.Order.StatusID != 2 {
			var order models.Order
			if err := s.db.First(&order, part.OrderID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return gorm.ErrRecordNotFound
				}
				return err
			}
			order.StatusID = 2
			result := s.db.Save(&order)
			if result.Error != nil {
				return result.Error
			}
		}
	}

	part.StatusID = statusID
	result := s.db.Save(&part)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
