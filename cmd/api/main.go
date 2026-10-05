package main

import (
	"context"
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
	logger.Info.Printf("PostgreSQL Connected")
	redisClient := cache.NewRedisClient(cfg.Redis)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}
	logger.Info.Printf("Redis Connected")

	// ===================== Repository =====================

	userRepo := postgres.NewUserRepository(db)
	genreRepo := postgres.NewGenreRepository(db)
	filmRepo := postgres.NewFilmRepository(db)
	seatRepo := postgres.NewSeatRepository(db)
	roomRepo := postgres.NewRoomRepository(db)
	ScheduleRepo := postgres.NewScheduleRepository(db)
	promoRepo := postgres.NewPromoRepository(db)
	scheduleSeatRepo := postgres.NewScheduleSeatRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)
	ticketRepo := postgres.NewTicketRepository(db)
	guestOrderRepo := postgres.NewGuestOrderRepository(db)

	//===================== Usecase =====================

	authUC := usecase.NewAuthUsecase(userRepo, redisClient, cfg.JWT)
	genreUC := usecase.NewGenreUsecase(genreRepo)
	filmUC := usecase.NewFilmUsecase(filmRepo, genreRepo)
	// seatUC := usecase.NewSeatUsecase(seatRepo)
	roomUC := usecase.NewRoomUsecase(roomRepo, seatRepo, db)
	scheduleSeatUC := usecase.NewScheduleSeatUsecase(scheduleSeatRepo)
	ScheduleUC := usecase.NewScheduleUsecase(ScheduleRepo, filmRepo, roomRepo)
	promoUC := usecase.NewPromoUsecase(promoRepo)
	bookingUC := usecase.NewBookingUsecase(db, ScheduleRepo, seatRepo, scheduleSeatRepo, transactionRepo, ticketRepo, guestOrderRepo, userRepo, cfg.Midtrans)

	//===================== handler =====================

	authHandler := handler.NewAuthHandler(authUC)
	genreHandler := handler.NewGenreHandler(genreUC)
	filmHandler := handler.NewFilmHandler(filmUC)
	// seatHandler := handler.NewSeatHandler(seatUC)
	roomHandler := handler.NewRoomHandler(roomUC)
	ScheduleHandler := handler.NewScheduleHandler(ScheduleUC)
	promoHandler := handler.NewPromoHandler(promoUC)
	homeHandler := handler.NewHomeHandler(authUC, filmUC, ScheduleUC, scheduleSeatUC, promoUC)
	bookingHandler := handler.NewBookingHandler(bookingUC)

	r := router.SetupRouter(
		authHandler,
		filmHandler,
		genreHandler,
		roomHandler,
		ScheduleHandler,
		promoHandler,
		homeHandler,
		bookingHandler,
		cfg.JWT.Secret,
		redisClient,
	)

	logger.Info.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
