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
	jwtSecret string,
	redisClient *redis.Client,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public film routes
		films := api.Group("/films")
		{
			films.GET("", filmHandler.GetAllWithDetails)
			films.GET("/:id", filmHandler.GetByID)
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
			// admin.GET("/Media", genreHandler.GetAll)
			// admin.GET("/genres", genreHandler.GetAll)
			// admin.GET("/genres", genreHandler.GetAll)
		}
	}

	return r
}
