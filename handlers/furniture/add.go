package furniture

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"io"
	"log/slog"
	"mis/storage/models"
	"mis/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type AddRequest struct {
	Name        string  `form:"name"`
	Price       uint    `form:"price"`
	Description *string `form:"description,omitempty"`
	Photo       *string `form:"photo,omitempty"`
	CategoryID  uint    `form:"categoryID"`
}

type AddResponse struct {
	ID uint `json:"id"`
}

func Add(log *slog.Logger, db FurnitureRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.furniture.Add"
		log := log.With(slog.String("fn", fn))

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Error("failed to parse multipart form", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("invalid request body"))
			return
		}

		var req AddRequest
		req.Name = r.FormValue("name")
		priceStr := r.FormValue("price")
		if priceStr != "" {
			price, err := strconv.ParseUint(priceStr, 10, 32)
			if err != nil {
				log.Info("invalid price format")
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, utils.NewErrorResponse("invalid price format"))
				return
			}
			req.Price = uint(price)
		}
		categoryIDStr := r.FormValue("categoryID")
		if categoryIDStr != "" {
			categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
			if err != nil {
				log.Info("invalid categoryID format")
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, utils.NewErrorResponse("invalid categoryID format"))
				return
			}
			req.CategoryID = uint(categoryID)
		}
		if desc := r.FormValue("description"); desc != "" {
			req.Description = &desc
		}

		if req.Name == "" || req.Price == 0 || req.CategoryID == 0 {
			log.Info("missing required fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("name, price, and categoryID are required"))
			return
		}

		file, handler, err := r.FormFile("photo")
		var photoPath *string
		if err == nil {
			defer file.Close()

			ext := filepath.Ext(handler.Filename)
			if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
				log.Warn("unsupported file format", slog.String("filename", handler.Filename))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, utils.NewErrorResponse("unsupported file format"))
				return
			}

			filename := uuid.New().String() + ext
			savePath := filepath.Join("/app/uploads", filename)

			if err := os.MkdirAll("/app/uploads", os.ModePerm); err != nil {
				log.Error("failed to create directory", slog.String("error", err.Error()))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, utils.NewErrorResponse("failed to create directory"))
				return
			}

			dst, err := os.Create(savePath)
			if err != nil {
				log.Error("failed to create file", slog.String("error", err.Error()))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, utils.NewErrorResponse("failed to save file"))
				return
			}
			defer dst.Close()

			if _, err := io.Copy(dst, file); err != nil {
				log.Error("failed to save file", slog.String("error", err.Error()))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, utils.NewErrorResponse("failed to save file"))
				return
			}

			photoPath = &filename
			log.Info("file uploaded successfully", slog.String("filename", filename))
		} else if err != http.ErrMissingFile {
			log.Error("failed to get file", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("failed to process file"))
			return
		}

		resultID, err := db.AddFurniture(models.Furniture{
			Name:        req.Name,
			CategoryID:  req.CategoryID,
			Price:       req.Price,
			Description: req.Description,
			Photo:       photoPath,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
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
		render.JSON(w, r, utils.NewSuccessResponse(AddResponse{ID: resultID}))
	}
}
