package orders

import (
	"github.com/go-chi/render"
	"log/slog"
	"mis/storage/models"
	"mis/utils"
	"net/http"
	"time"
)

type PartOfOrder struct {
	FurnitureID int `json:"furnitureID"`
	Count       int `json:"count"`
}

type AddRequest struct {
	CustomerFirstName string        `json:"customerFirstName"`
	CustomerLastName  string        `json:"customerLastName"`
	CustomerPhone     string        `json:"customerPhone"`
	Furniture         []PartOfOrder `json:"furniture"`
}

func Add(log *slog.Logger, db OrdersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.Add"
		log := log.With(slog.String("fn", fn))

		var req AddRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "invalid request body"})
			return
		}
		if req.CustomerFirstName == "" || req.CustomerLastName == "" || req.CustomerPhone == "" || len(req.Furniture) == 0 {
			log.Info("missing required fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("customerFirstName, customerLastName, customerPhone and furniture are required"))
			return
		}
		sum := 0
		var partsOfOrder []models.PartOfOrder
		for _, part := range req.Furniture {
			furniture, err := db.GetFurnitureByID(uint(part.FurnitureID))
			if err != nil {
				log.Error("failed to get furniture by ID")
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, utils.NewErrorResponse("failed to get furniture by ID"))
				return
			}
			sum += int(furniture.Price) * part.Count
			for i := 0; i < part.Count; i++ {
				partsOfOrder = append(partsOfOrder, models.PartOfOrder{
					FurnitureID: uint(part.FurnitureID),
					StatusID:    1,
				})
			}
		}

		orderID, err := db.AddOrder(models.Order{
			CreatedAt:         time.Now(),
			CustomerFirstName: req.CustomerFirstName,
			CustomerLastName:  req.CustomerLastName,
			CustomerPhone:     req.CustomerPhone,
			Sum:               uint(sum),
			StatusID:          1,
			PartsOfOrder:      partsOfOrder,
		})
		if err != nil {
			log.Error("failed to add order", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("failed to add order"))
			return
		}
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, utils.NewSuccessResponse(map[string]uint{"orderID": orderID}))
	}
}
