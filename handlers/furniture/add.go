package furniture

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
	"mis/storage/models"
	"mis/utils"
	"net/http"
)

type AddRequest struct {
	Name        string  `json:"name"`
	Price       uint    `json:"price"`
	Description *string `json:"description,omitempty"`
	Photo       *string `json:"photo,omitempty"`
	CategoryID  uint    `json:"categoryID"`
}

type AddResponse struct {
	ID uint `json:"id"`
}

func Add(log *slog.Logger, db FurnitureRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.furniture.Add"
		log := log.With(slog.String("fn", fn))

		var req AddRequest
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

		resultId, err := db.AddFurniture(models.Furniture{
			Name:        req.Name,
			CategoryID:  req.CategoryID,
			Price:       req.Price,
			Description: req.Description,
			Photo:       req.Photo,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				fmt.Println("true")
				log.Info("furniture with the given name already exists")
				render.Status(r, http.StatusConflict)
				render.JSON(w, r, utils.NewErrorResponse("furniture with the given name already exists"))
				return
			}
			log.Error("failed to add furniture", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("internal server error"))
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, utils.NewSuccessResponse(AddResponse{ID: resultId}))
	}
}
