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

func SetupRouterRoutes(routerGroup *gin.RouterGroup, handler *handler.RouterHandler) {
	routers := routerGroup.Group("/routers")
	{
		routers.POST("", handler.Create)
		routers.GET("", handler.GetAll)
		routers.GET("/:id", handler.GetByID)
		routers.PUT("/:id", handler.Update)
		routers.DELETE("/:id", handler.Delete)
	}
}

func SetupServiceProfileRoutes(routerGroup *gin.RouterGroup, handler *handler.ServiceProfileHandler) {
	profiles := routerGroup.Group("/profiles")
	{
		profiles.POST("", handler.Create)
		profiles.GET("", handler.GetAll)
		profiles.GET("/:id", handler.GetByID)
		profiles.PUT("/:id", handler.Update)
		profiles.DELETE("/:id", handler.Delete)
	}
}

func SetupCustomerRoutes(router *gin.RouterGroup, handler *handler.CustomerHandler) {
	customerRoutes := router.Group("/customers")
	{
		customerRoutes.POST("", handler.CreateCustomer)
		customerRoutes.GET("", handler.GetAllCustomers)
		customerRoutes.GET("/:id", handler.GetCustomer)
		customerRoutes.PUT("/:id", handler.UpdateCustomer)
		customerRoutes.POST("/:id/suspend", handler.SuspendCustomer)
		customerRoutes.POST("/:id/activate", handler.ActivateCustomer)
		customerRoutes.DELETE("/:id", handler.DeleteCustomer)
	}
}

func SetupVoucherRoutes(routerGroup *gin.RouterGroup, handler *handler.VoucherHandler) {
	vouchers := routerGroup.Group("/vouchers")
	{
		vouchers.POST("/generate", handler.Generate)
		vouchers.GET("", handler.GetAll)
		vouchers.GET("/:id", handler.GetByID)
		vouchers.DELETE("/:id", handler.Delete)
	}
}