package furniture

import (
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"mis/storage/models"
	"mis/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type UpdateRequest struct {
	ID          uint    `form:"id"`
	Name        string  `form:"name"`
	Description *string `form:"description,omitempty"`
	Photo       *string `form:"photo,omitempty"`
	Price       uint    `form:"price"`
	CategoryID  uint    `form:"categoryID"`
}

func Update(log *slog.Logger, repo FurnitureRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.furniture.Update"
		log := log.With(slog.String("fn", fn))

		// Ограничение размера формы (10MB)
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Error("failed to parse multipart form", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("invalid request body"))
			return
		}

		// Извлечение данных из формы
		var req UpdateRequest
		idStr := r.FormValue("id")
		if idStr != "" {
			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				log.Info("invalid id format")
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, utils.NewErrorResponse("invalid id format"))
				return
			}
			req.ID = uint(id)
		}
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

		// Проверка обязательных полей
		if req.ID == 0 || req.Name == "" || req.Price == 0 || req.CategoryID == 0 {
			log.Info("missing required fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id, name, price, and categoryID are required"))
			return
		}

		// Обработка файла (если загружен)
		file, handler, err := r.FormFile("photo")
		var photoPath *string
		if err == nil {
			defer file.Close()

			// Проверка расширения
			ext := filepath.Ext(handler.Filename)
			if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
				log.Warn("unsupported file format", slog.String("filename", handler.Filename))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, utils.NewErrorResponse("unsupported file format"))
				return
			}

			// Генерация уникального имени файла
			filename := uuid.New().String() + ext
			savePath := filepath.Join("/app/uploads", filename)

			// Создание папки uploads
			if err := os.MkdirAll("/app/uploads", os.ModePerm); err != nil {
				log.Error("failed to create directory", slog.String("error", err.Error()))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, utils.NewErrorResponse("failed to create directory"))
				return
			}

			// Сохранение файла
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

		// Если фото не загружено, оставляем текущее (или очищаем, если нужно)
		// Здесь предполагается, что Photo в базе не обновляется, если файл не предоставлен
		furniture := models.Furniture{
			ID:          req.ID,
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			CategoryID:  req.CategoryID,
		}
		if photoPath != nil {
			furniture.Photo = photoPath
		}

		// Обновление в базе данных
		err = repo.UpdateFurniture(furniture)
		if err != nil {
			log.Error("failed to update furniture", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("failed to update furniture"))
			return
		}

		render.Status(r, http.StatusOK)
		//render.JSON(w, r, utils.NewSuccessResponse())
	}
}
