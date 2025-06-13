package furniture

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"mis/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func GetPhoto(log *slog.Logger, repo FurnitureRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.furniture.GetPhoto"
		log := log.With(slog.String("fn", fn))

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			log.Info("invalid furniture ID", slog.String("id", idStr))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("invalid furniture ID"))
			return
		}

		furniture, err := repo.GetFurnitureByID(uint(id))
		if err != nil {
			log.Error("failed to get furniture", slog.String("error", err.Error()))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, utils.NewErrorResponse("furniture not found"))
			return
		}

		if furniture.Photo == nil || *furniture.Photo == "" {
			log.Info("no photo for furniture", slog.Uint64("id", uint64(id)))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, utils.NewErrorResponse("no photo available"))
			return
		}

		filePath := filepath.Join("/app/uploads", *furniture.Photo)
		file, err := os.Open(filePath)
		if err != nil {
			log.Error("failed to open file", slog.String("path", filePath), slog.String("error", err.Error()))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, utils.NewErrorResponse("photo not found"))
			return
		}
		defer file.Close()

		ext := filepath.Ext(*furniture.Photo)
		var contentType string
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		default:
			log.Warn("unsupported file extension", slog.String("ext", ext))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("unsupported file format"))
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Cache-Control", "public, max-age=86400")

		if _, err := io.Copy(w, file); err != nil {
			log.Error("failed to send file", slog.String("error", err.Error()))
			return
		}

		log.Info("photo sent successfully", slog.Uint64("id", uint64(id)), slog.String("filename", *furniture.Photo))
	}
}
