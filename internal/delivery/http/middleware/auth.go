package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	jwtPkg "github.com/rafli/boocins/pkg/jwt"
	"github.com/rafli/boocins/pkg/response"
)

func AuthMiddleware(jwtSecret string, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error("authorization header required"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error("invalid authorization format"))
			return
		}

		tokenString := parts[1]

		blacklistKey := fmt.Sprintf("blacklist:%s", tokenString)
		if redisClient.Exists(c.Request.Context(), blacklistKey).Val() > 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error("token has been revoked"))
			return
		}

		claims, err := jwtPkg.ParseToken(tokenString, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error("invalid or expired token"))
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("token", tokenString)
		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error("role not found"))
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error("invalid role"))
			return
		}

		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, response.Error("access denied"))
	}
}
