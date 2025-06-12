package orders

import (
	"github.com/go-chi/render"
	"log/slog"
	"mis/utils"
	"net/http"
)

type PatchPartRequest struct {
	PartID   uint `json:"partID"`
	StatusID uint `json:"statusID"`
}

func PatchPart(log *slog.Logger, db OrdersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.PatchPart"
		log := log.With(slog.String("fn", fn))

		var req PatchPartRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "invalid request body"})
			return
		}
		if req.PartID == 0 || req.StatusID == 0 {
			log.Info("missing required fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "partID and statusID are required"})
			return
		}
		if err := db.UpdatePartOfOrderStatus(req.PartID, req.StatusID); err != nil {
			log.Error("failed to update part of order status", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("failed to update part of order status"))
			return
		}
	}
}
