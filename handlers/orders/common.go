package orders

import "mis/storage/models"

type OrdersRepo interface {
	GetOrders() ([]models.Order, error)
	GetOrderByID(orderID uint) (models.Order, error)
	AddOrder(order models.Order) (uint, error)
	UpdateOrder(order models.Order) error
	GetFurnitureByID(id uint) (models.Furniture, error)
	UpdateOrderStatus(orderID uint, statusID uint) error
	UpdatePartOfOrderStatus(partID uint, statusID uint) error
}
