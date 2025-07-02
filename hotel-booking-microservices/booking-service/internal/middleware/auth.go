package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"booking-service/internal/services"
)

// AuthMiddleware verifica JWT válido
func AuthMiddleware(bookingService *services.BookingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token de autorización requerido",
			})
			c.Abort()
			return
		}

		// Verificar formato Bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Formato de token inválido",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validar token y obtener datos del usuario
		userID, role, err := bookingService.ValidateJWTTokenWithRole(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token inválido",
			})
			c.Abort()
			return
		}

		// Guardar datos en contexto
		c.Set("userID", userID)
		c.Set("userRole", role)
		c.Next()
	}
}

// AdminMiddleware verifica que el usuario sea admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token requerido",
			})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: se requieren permisos de administrador",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromContext obtiene datos del usuario del contexto
func GetUserFromContext(c *gin.Context) (int, string, bool) {
	userID, exists1 := c.Get("userID")
	userRole, exists2 := c.Get("userRole")
	
	if !exists1 || !exists2 {
		return 0, "", false
	}
	
	id, ok1 := userID.(int)
	role, ok2 := userRole.(string)
	
	return id, role, ok1 && ok2
}