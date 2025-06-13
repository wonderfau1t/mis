package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"log"
	"log/slog"
	"mis/handlers/categories"
	"mis/handlers/furniture"
	"mis/handlers/orders"
	"mis/handlers/statuses"
	"mis/storage"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello, World!")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	db, err := storage.SetupStorage()
	if err != nil {
		log.Fatalf("Faield to set storage")
	}
	_ = logger

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"message": "Hello, World!"})
	})

	router.Get("/statuses", statuses.List(logger, db))
	router.Get("/categories", categories.List(logger, db))

	router.Get("/furniture", furniture.List(logger, db))
	router.Get("/furniture/{id}", furniture.Get(logger, db))
	router.Post("/furniture", furniture.Add(logger, db))
	router.Put("/furniture", furniture.Update(logger, db))

	router.Get("/orders", orders.List(logger, db))
	router.Post("/orders", orders.Add(logger, db))
	router.Put("/orders", orders.Update(logger, db))

	router.Patch("/orders", orders.Patch(logger, db))
	router.Patch("/orders/part", orders.PatchPart(logger, db))

	err = http.ListenAndServe("0.0.0.0:8000", router)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
