package orders

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"log/slog"
	"mis/utils"
	"net/http"
	"time"
)

type PartOfOrderDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	StatusID uint   `json:"statusID"`
}

type OrderDTO struct {
	ID                uint             `json:"id"`
	CreatedAt         time.Time        `json:"createdAt"`
	CustomerFirstName string           `json:"customerFirstName"`
	CustomerLastName  string           `json:"customerLastName"`
	CustomerPhone     string           `json:"customerPhone"`
	Sum               uint             `json:"sum"`
	StatusID          uint             `json:"statusID"`
	Furniture         []PartOfOrderDTO `json:"furniture"`
}

type ListResponse struct {
	TotalCount int        `json:"totalCount"`
	Orders     []OrderDTO `json:"orders"`
}

func List(log *slog.Logger, repo OrdersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.List"
		log := log.With(slog.String("fn", fn))

		orders, err := repo.GetOrders()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Info("no orders found")
				render.Status(r, http.StatusOK)
				render.JSON(w, r, utils.NewSuccessResponse(ListResponse{
					TotalCount: 0,
					Orders:     nil,
				}))
				return
			}
			log.Error("failed to get list of orders", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			return
		}

		ordersDTO := make([]OrderDTO, 0, len(orders))
		for i, order := range orders {
			fmt.Println(order.ID)
			ordersDTO = append(ordersDTO, OrderDTO{
				ID:                order.ID,
				CreatedAt:         order.CreatedAt,
				CustomerFirstName: order.CustomerFirstName,
				CustomerLastName:  order.CustomerLastName,
				CustomerPhone:     order.CustomerPhone,
				Sum:               order.Sum,
				StatusID:          order.StatusID,
			})
			ordersDTO[i].Furniture = make([]PartOfOrderDTO, len(order.PartsOfOrder))
			for j, part := range order.PartsOfOrder {
				ordersDTO[i].Furniture[j] = PartOfOrderDTO{
					ID:       part.ID,
					Name:     part.Furniture.Name,
					StatusID: part.StatusID,
				}
			}
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(ListResponse{
			TotalCount: len(orders),
			Orders:     ordersDTO,
		}))
	}
}
