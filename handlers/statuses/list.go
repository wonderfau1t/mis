package statuses

import (
	"github.com/go-chi/render"
	"log/slog"
	"mis/storage/models"
	"mis/utils"
	"net/http"
)

type ListResponse struct {
	Statuses []models.Status `json:"statuses"`
}

type CategoriesRepo interface {
	GetStatuses() ([]models.Status, error)
}

func List(log *slog.Logger, repo CategoriesRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.categories.List"
		log := log.With(slog.String("fn", fn))

		statuses, err := repo.GetStatuses()
		if err != nil {
			log.Error("failed to get list of categories")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(ListResponse{Statuses: statuses}))
	}
}
