package middleware

import (
	"net/http"

	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/utils"
	"github.com/gin-gonic/gin"
)

// La funzione AuthMiddleware protegge le rotte private verificando la presenza e la validità del token JWT
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		// tenta di estrarre il token di accesso
		token, err := utils.GetAccessToken(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// valida la firma digitale, l'integrità e la data di scadenza del token JWT
		claims, err := utils.ValidateToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)

		c.Next()

	}
}
