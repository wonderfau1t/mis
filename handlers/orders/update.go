package orders

import (
	"github.com/go-chi/render"
	"log/slog"
	"mis/storage/models"
	"mis/utils"
	"net/http"
	"time"
)

type UpdateRequest struct {
	ID uint `json:"id"`
	AddRequest
}

func Update(log *slog.Logger, db OrdersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.Update"
		log := log.With(slog.String("fn", fn))

		var req UpdateRequest
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

		err := db.UpdateOrder(models.Order{
			ID:                req.ID,
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
		render.Status(r, http.StatusOK)
		//render.JSON(w, r, utils.NewSuccessResponse(map[string]uint{"orderID": orderID}))
	}
}
