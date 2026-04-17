package middleware

import (
	"net/http"

	"super-indo-api/pkg/infrastructure/jwt"
	"super-indo-api/pkg/common"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSvc *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Error(http.StatusUnauthorized, "token tidak ditemukan"))
			return
		}

		claims, err := jwtSvc.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Error(http.StatusUnauthorized, "token tidak valid atau sudah expired"))
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
