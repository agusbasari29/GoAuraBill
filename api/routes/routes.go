package routes
import (
	"github.com/gin-gonic/gin"
	"github.com/agusbasari29/GoAuraBill/internal/handler" // Ganti dengan path modul Anda
	"github.com/agusbasari29/GoAuraBill/internal/middleware" // Ganti dengan path modul Anda
)
func SetupAuthRoutes(router *gin.Engine, authHandler *handler.AuthHandler, jwtSecret string) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}
	// Contoh rute yang dilindungi
	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protectedRoutes.GET("/profile", func(c *gin.Context) {
			userID := middleware.GetUserID(c)
			role := middleware.GetUserRole(c)
			c.JSON(200, gin.H{"message": "Welcome to your profile!", "user_id": userID, "role": role})
		})
	}
}