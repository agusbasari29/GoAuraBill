package main
import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/agusbasari29/GoAuraBill/config"
	"github.com/agusbasari29/GoAuraBill/internal/handler"
	"github.com/agusbasari29/GoAuraBill/internal/middleware"
	"github.com/agusbasari29/GoAuraBill/internal/model"
	"github.com/agusbasari29/GoAuraBill/internal/repository"
	"github.com/agusbasari29/GoAuraBill/internal/service"
	"github.com/agusbasari29/GoAuraBill/routes" // Import paket routes
)
var DB *gorm.DB
func main() {
	// 1. Muat Konfigurasi
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}
	// 2. Hubungkan ke Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connection established")
	// 3. Jalankan Migrasi Otomatis
	log.Println("Running Migrations")
	err = DB.AutoMigrate(
		&model.User{},
		&model.Router{},
		&model.ServiceProfile{},
		&model.Voucher{},
		&model.Transaction{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("Database Migrated Successfully")
	// 4. Inisialisasi Dependensi

	// Auth
	authRepo := repository.NewAuthRepository(DB)
	authService := service.NewAuthService(authRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)
	// Router
	routerRepo := repository.NewRouterRepository(DB)
	routerService := service.NewRouterService(routerRepo)
	routerHandler := handler.NewRouterHandler(routerService)
	// ServiceProfile
	profileRepo := repository.NewServiceProfileRepository(DB)
	profileService := service.NewServiceProfileService(profileRepo)
	profileHandler := handler.NewServiceProfileHandler(profileService)

	customerRepo := repository.NewCustomerRepository(DB)
	customerService := service.NewCustomerService(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerService)

	
	transactionRepo := repository.NewTransactionRepository(DB)
	transactionService := service.NewTransactionService(transactionRepo, customerRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	
	voucherRepo := repository.NewVoucherRepository(DB)
	voucherService := service.NewVoucherService(voucherRepo, customerRepo, transactionRepo, DB)
	voucherHandler := handler.NewVoucherHandler(voucherService)

	paymentService := service.NewPaymentService(cfg, transactionRepo, customerRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	reportRepo := repository.NewReportRepository(DB)
	reportService := service.NewReportService(reportRepo)
	reportHandler := handler.NewReportHandler(reportService)

	// 5. Siapkan Server Gin
	router := gin.Default()
	// 6. Setup Rute
	routes.SetupAuthRoutes(router, authHandler, cfg.JWTSecret) // Panggil fungsi setup rute
	// Grup rute yang dilindungi
	apiRoutes := router.Group("/api")
	apiRoutes.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		routes.SetupRouterRoutes(apiRoutes, routerHandler)
		routes.SetupServiceProfileRoutes(apiRoutes, profileHandler)
		routes.SetupCustomerRoutes(apiRoutes, customerHandler)
		routes.SetupVoucherRoutes(apiRoutes, voucherHandler)
		routes.SetupTransactionRoutes(apiRoutes, transactionHandler)
		routes.SetupPaymentRoutes(router, apiRoutes, paymentHandler)
		routes.SetupReportRoutes(apiRoutes, reportHandler)
	}
	// 7. Jalankan Server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}	