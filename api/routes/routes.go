package routes

import (
	"github.com/agusbasari29/GoAuraBill/internal/handler"    // Ganti dengan path modul Anda
	"github.com/agusbasari29/GoAuraBill/internal/middleware" // Ganti dengan path modul Anda
	"github.com/gin-gonic/gin"
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
		vouchers.POST("/generate", handler.Generate) // Admin only
		vouchers.GET("", handler.GetAll)             // Admin only
		vouchers.GET("/:id", handler.GetByID)        // Admin only
		vouchers.DELETE("/:id", handler.Delete)      // Admin only
		vouchers.POST("/redeem", handler.Redeem)     // Customer
	}
}

func SetupTransactionRoutes(router *gin.RouterGroup, handler *handler.TransactionHandler) {
	txns := router.Group("/transactions")
	{
		txns.POST("", handler.Create)
		txns.GET("/:id", handler.GetByID)
		txns.GET("/customer/:customer_id", handler.GetCustomerTransactions)
		txns.POST("/:id/process", handler.ProcessPayment)
		txns.POST("/:id/cancel", handler.CancelTransaction)
	}
}

func SetupPaymentRoutes(router *gin.Engine, apiGroup *gin.RouterGroup, handler *handler.PaymentHandler) {
	// Endpoint untuk membuat charge, memerlukan autentikasi
	apiGroup.POST("/payments/charge/:transaction_id", handler.CreateCharge)
	
	// Endpoint untuk callback dari Tripay, TIDAK memerlukan autentikasi
	router.POST("/api/payments/tripay-callback", handler.HandleNotification)
}

func SetupReportRoutes(routerGroup *gin.RouterGroup, handler *handler.ReportHandler) {
	reports := routerGroup.Group("/reports")
	{
		reports.GET("/revenue", handler.GetRevenueReport) // GET /api/reports/revenue?period=daily/monthly
		reports.GET("/summary", handler.GetSummaryReport) // GET /api/reports/summary
	}
}