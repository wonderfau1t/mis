package furniture

import (
	"github.com/go-chi/render"
	"log/slog"
	"mis/utils"
	"net/http"
	"strconv"
)

type GetResponse struct {
	Furniture FurnitureDTO `json:"furniture"`
}

func Get(log *slog.Logger, db FurnitureRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.furniture.Get"
		log := log.With(slog.String("fn", fn))
		id := r.URL.Query().Get("id")
		intID, _ := strconv.Atoi(id)
		furniture, err := db.GetFurnitureByID(uint(intID))
		if err != nil {
			log.Error("failed to get furniture by ID", slog.Any("error", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		furnitureDTO := FurnitureDTO{
			ID:           furniture.ID,
			Name:         furniture.Name,
			CategoryName: furniture.Category.Name,
			Price:        furniture.Price,
			Photo:        furniture.Photo,
			Description:  furniture.Description,
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(GetResponse{Furniture: furnitureDTO}))
	}
}
