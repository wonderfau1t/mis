package furniture

import (
	"github.com/go-chi/render"
	"log/slog"
	"mis/storage/models"
	"mis/utils"
	"net/http"
)

type ListResponse struct {
	TotalCount uint               `json:"totalCount"`
	Furniture  []models.Furniture `json:"furniture"`
}

func List(log *slog.Logger, repo FurnitureRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.furniture.List"
		log := log.With(slog.String("fn", fn))

		furniture, err := repo.GetFurniture()
		if err != nil {
			log.Error("failed to get list of furniture")
			render.Status(r, http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(ListResponse{
			TotalCount: uint(len(furniture)),
			Furniture:  furniture,
		}))
	}
}
