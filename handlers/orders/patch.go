package orders

import (
	"github.com/go-chi/render"
	"log/slog"
	"mis/utils"
	"net/http"
)

type PatchRequest struct {
	OrderID  uint `json:"orderID"`
	StatusID uint `json:"statusID"`
}

func Patch(log *slog.Logger, db OrdersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.Patch"
		log := log.With(slog.String("fn", fn))

		var req PatchRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "invalid request body"})
			return
		}

		if req.OrderID == 0 || req.StatusID == 0 {
			log.Info("missing required fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("orderID and statusID are required"))
			return
		}

		if err := db.UpdateOrderStatus(req.OrderID, req.StatusID); err != nil {
			log.Error("failed to update order status", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("failed to update order status"))
			return
		}
	}
}
