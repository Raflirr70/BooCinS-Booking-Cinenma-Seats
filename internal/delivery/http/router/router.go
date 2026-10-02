package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rafli/boocins/internal/delivery/http/handler"
	"github.com/rafli/boocins/internal/delivery/http/middleware"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	filmHandler *handler.FilmHandler,
	genreHandler *handler.GenreHandler,
	roomHandler *handler.RoomHandler,
	scheduleHandler *handler.ScheduleHandler,
	promoHandler *handler.PromoHandler,
	homeHandler *handler.HomeHandler,
	jwtSecret string,
	redisClient *redis.Client,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/home", homeHandler.GetHome)
		api.GET("/schedules/:id/seats", homeHandler.GetScheduleMapSeats)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public film routes
		films := api.Group("/films")
		{
			films.GET("", homeHandler.GetFilms)
			films.GET("/:id", homeHandler.GetDetailFilms)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret, redisClient))
		{
			protected.POST("/auth/logout", authHandler.Logout)
			protected.GET("/auth/profile", authHandler.GetProfile)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(jwtSecret, redisClient))
		admin.Use(middleware.RoleMiddleware("super_admin", "admin"))
		{
			admin.GET("/films", filmHandler.GetAll)
			admin.POST("/films", filmHandler.Create)
			admin.GET("/films/:id", filmHandler.GetByID)
			admin.PUT("/films/:id", filmHandler.Update)
			admin.DELETE("/films/:id", filmHandler.Delete)

			admin.GET("/genres", genreHandler.GetAll)
			admin.POST("/genres", genreHandler.Create)
			admin.GET("/genres/:id", genreHandler.GetByID)
			admin.PUT("/genres/:id", genreHandler.Update)
			admin.DELETE("/genres/:id", genreHandler.Delete)

			admin.POST("/room", roomHandler.Create)
			admin.GET("/room", roomHandler.GetAll)
			admin.PUT("/room/:id", roomHandler.Update)
			admin.DELETE("/room/:id", roomHandler.Delete)
			admin.GET("/room/:id", roomHandler.GetWithDetail)
			// admin.GET("/Media", genreHandler.GetAll)
			// admin.GET("/genres", genreHandler.GetAll)
			// admin.GET("/genres", genreHandler.GetAll)

			admin.GET("/schedule", scheduleHandler.GetAll)
			admin.POST("/schedule", scheduleHandler.Create)
			admin.PUT("/schedule/:id", scheduleHandler.Update)
			admin.DELETE("/schedule/:id", scheduleHandler.Delete)
			admin.GET("/schedule/:id/film", scheduleHandler.GetByFilm)
			admin.GET("/schedule/:id/room", scheduleHandler.GetByRoom)

			api.POST("/promo", promoHandler.Create)
			api.GET("/promo", promoHandler.GetCurrent)
			api.GET("/promo/:id", promoHandler.GetByID)
			api.PUT("/promo/:id", promoHandler.Update)
			api.DELETE("/promo/:id", promoHandler.Delete)
		}
	}

	return r
}
