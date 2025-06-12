package categories

import (
	"github.com/go-chi/render"
	"log/slog"
	"mis/storage/models"
	"mis/utils"
	"net/http"
)

type ListResponse struct {
	Categories []models.Category `json:"categories"`
}

type CategoriesRepo interface {
	GetCategories() ([]models.Category, error)
}

func List(log *slog.Logger, repo CategoriesRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.categories.List"
		log := log.With(slog.String("fn", fn))

		categories, err := repo.GetCategories()
		if err != nil {
			log.Error("failed to get list of categories")
			render.Status(r, http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(ListResponse{
			Categories: categories,
		}))
	}
}
