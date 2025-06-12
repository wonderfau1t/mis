package furniture

import (
	"github.com/go-chi/render"
	"log/slog"
	"mis/storage/models"
	"mis/utils"
	"net/http"
)

type UpdateRequest struct {
	Id          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Photo       *string `json:"photo,omitempty"`
	Price       uint    `json:"price"`
	CategoryID  uint    `json:"categoryID"`
}

func Update(log *slog.Logger, repo FurnitureRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.furniture.Update"
		log := slog.With(slog.String("fn", fn))

		var req UpdateRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("invalid request body"))
			return
		}

		if req.Name == "" || req.Price == 0 || req.CategoryID == 0 {
			log.Info("missing required fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("name, price, and categoryId are required"))
			return
		}

		err := repo.UpdateFurniture(models.Furniture{
			ID:          req.Id,
			Name:        req.Name,
			Description: req.Description,
			Photo:       req.Photo,
			Price:       req.Price,
			CategoryID:  req.CategoryID,
		})
		if err != nil {
			log.Error("failed to update furniture", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("failed to update furniture"))
			return
		}

		render.Status(r, http.StatusOK)
	}
}
