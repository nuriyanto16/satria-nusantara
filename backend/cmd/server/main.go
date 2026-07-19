package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"satria-nusantara/backend/internal/auth"
	"satria-nusantara/backend/internal/content"
	"satria-nusantara/backend/internal/event"
	"satria-nusantara/backend/internal/finance"
	"satria-nusantara/backend/internal/kebugaran"
	"satria-nusantara/backend/internal/organization"
	"satria-nusantara/backend/internal/training"

	"satria-nusantara/backend/pkg/config"
	"satria-nusantara/backend/pkg/database"
	"satria-nusantara/backend/pkg/response"
)

func main() {
	// 1. Load config
	cfg := config.Load()

	// 2. Init DB & Redis (will fatal if failed)
	db := database.NewPostgres(cfg)
	defer db.Close()
	rdb := database.NewRedis(cfg)
	defer rdb.Close()

	// 3. Initialize Repositories
	authRepo := auth.NewRepository(db)
	orgRepo := organization.NewRepository(db)
	trainRepo := training.NewRepository(db)

	// 4. Initialize Services
	authSvc := auth.NewService(authRepo, cfg.JWTSecret)
	orgSvc := organization.NewService(orgRepo)
	trainSvc := training.NewService(trainRepo, orgRepo)

	// 5. Initialize Handlers
	authHandler := auth.NewHandler(authSvc)
	orgHandler := organization.NewHandler(orgSvc)
	trainHandler := training.NewHandler(trainSvc)
	eventHandler := event.NewHandler()
	financeHandler := finance.NewHandler()
	contentHandler := content.NewHandler()
	kebugaranHandler := kebugaran.NewHandler()

	// 6. Setup Router (Chi)
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Serve static uploads
	fs := http.FileServer(http.Dir("./uploads"))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fs))

	// 7. Mount API Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			response.JSON(w, http.StatusOK, map[string]string{
				"status": "ok",
				"time":   time.Now().Format(time.RFC3339),
				"mode":   "modular-monolith",
			})
		})

		r.Route("/auth", authHandler.Routes(cfg.JWTSecret))
		r.Route("/organization", orgHandler.Routes(cfg.JWTSecret))
		r.Route("/training", trainHandler.Routes(cfg.JWTSecret))
		r.Route("/event", eventHandler.Routes())
		r.Route("/finance", financeHandler.Routes())
		r.Route("/content", contentHandler.Routes())
		r.Route("/kebugaran", kebugaranHandler.Routes())
	})

	// 8. Start Server
	log.Printf("🚀 Server started on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("server: failed to listen: %v", err)
	}
}
