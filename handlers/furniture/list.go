package furniture

import (
	"fmt"
	"github.com/go-chi/render"
	"log/slog"
	"mis/utils"
	"net/http"
)

type FurnitureDTO struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	CategoryName string  `json:"categoryName"`
	Price        uint    `json:"price"`
	Photo        string  `json:"photo,omitempty"`
	Description  *string `json:"description"`
}

type ListResponse struct {
	TotalCount uint           `json:"totalCount"`
	Furniture  []FurnitureDTO `json:"furniture"`
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

		furnitureDTOs := make([]FurnitureDTO, 0, len(furniture))

		for _, item := range furniture {
			furnitureDTOs = append(furnitureDTOs, FurnitureDTO{
				ID:           item.ID,
				Name:         item.Name,
				CategoryName: item.Category.Name,
				Price:        item.Price,
				Photo:        fmt.Sprintf("/furniture/%d/photo", item.ID),
				Description:  item.Description,
			})
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(ListResponse{
			TotalCount: uint(len(furniture)),
			Furniture:  furnitureDTOs,
		}))
	}
}
