package orders

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type FurnitureDTOtoUpdate struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type OrderDTOtoUpdate struct {
	ID                uint                   `json:"id"`
	CreatedAt         time.Time              `json:"createdAt"`
	CustomerFirstName string                 `json:"customerFirstName"`
	CustomerLastName  string                 `json:"customerLastName"`
	CustomerPhone     string                 `json:"customerPhone"`
	Sum               uint                   `json:"sum"`
	StatusID          uint                   `json:"statusID"`
	Furniture         []FurnitureDTOtoUpdate `json:"furniture"`
}

type GetResponse struct {
	Order OrderDTOtoUpdate `json:"order"`
}

func Get(log *slog.Logger, db OrdersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.Get"
		log := log.With(slog.String("fn", fn))

		id := chi.URLParam(r, "id")
		intID, _ := strconv.Atoi(id)
		order, err := db.GetOrderByID(uint(intID))
		if err != nil {
			log.Error("failed to get order by ID", slog.Any("error", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		orderDTO := OrderDTOtoUpdate{
			ID:                order.ID,
			CreatedAt:         order.CreatedAt,
			CustomerFirstName: order.CustomerFirstName,
			CustomerLastName:  order.CustomerLastName,
			CustomerPhone:     order.CustomerPhone,
			Sum:               order.Sum,
			StatusID:          order.StatusID,
		}
		var furnitureCount map[uint]int
		for _, part := range order.PartsOfOrder {
			if furnitureCount == nil {
				furnitureCount = make(map[uint]int)
			}
			furnitureCount[part.FurnitureID] += 1
		}
		for furnitureID, count := range furnitureCount {
			furniture, err := db.GetFurnitureByID(furnitureID)
			if err != nil {
				log.Error("failed to get furniture by ID", slog.Any("error", err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			orderDTO.Furniture = append(orderDTO.Furniture, FurnitureDTOtoUpdate{
				ID:    furniture.ID,
				Name:  furniture.Name,
				Count: count,
			})
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, GetResponse{Order: orderDTO})
	}
}
