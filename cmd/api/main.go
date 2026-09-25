package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/rafli/boocins/config"
	"github.com/rafli/boocins/internal/delivery/http/handler"
	"github.com/rafli/boocins/internal/delivery/http/router"
	"github.com/rafli/boocins/internal/infrastructure/cache"
	"github.com/rafli/boocins/internal/infrastructure/database"
	"github.com/rafli/boocins/internal/repository/postgres"
	"github.com/rafli/boocins/internal/usecase"
	"github.com/rafli/boocins/pkg/logger"
)

func main() {
	logger.Init()

	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found, using system env")
	}

	cfg := config.LoadConfig()

	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	redisClient := cache.NewRedisClient(cfg.Redis)

	userRepo := postgres.NewUserRepository(db)
	authUC := usecase.NewAuthUsecase(userRepo, redisClient, cfg.JWT)
	authHandler := handler.NewAuthHandler(authUC)

	filmRepo := postgres.NewFilmRepository(db)
	filmUC := usecase.NewFilmUsecase(filmRepo)
	filmHandler := handler.NewFilmHandler(filmUC)

	genreRepo := postgres.NewGenreRepository(db)
	genreUC := usecase.NewGenreUsecase(genreRepo)
	genreHandler := handler.NewGenreHandler(genreUC)

	r := router.SetupRouter(
		authHandler,
		filmHandler,
		genreHandler,
		cfg.JWT.Secret,
		redisClient,
	)

	logger.Info.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
