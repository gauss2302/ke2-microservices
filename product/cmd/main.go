// product/cmd/main.go
package main

import (
	"github.com/gauss2302/testcommm/product/internal/domain/entity"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gauss2302/testcommm/product/internal/handler"
	authMiddleware "github.com/gauss2302/testcommm/product/internal/middleware"
	"github.com/gauss2302/testcommm/product/internal/repository"
	"github.com/gauss2302/testcommm/product/internal/service"
	pb "github.com/gauss2302/testcommm/product/proto/auth"
)

func main() {
	// Database connection with retry
	var db *gorm.DB
	var err error

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database, attempt %d/%d: %v", i+1, maxRetries, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	log.Printf("Successfully connected to database")

	// Auto migrate
	if err := db.AutoMigrate(&entity.Product{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Connect to auth service with retry
	var authConn *grpc.ClientConn
	for i := 0; i < maxRetries; i++ {
		authConn, err = grpc.Dial(
			os.Getenv("AUTH_SERVICE_ADDR"),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to auth service, attempt %d/%d: %v", i+1, maxRetries, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to auth service after %d attempts: %v", maxRetries, err)
	}
	defer authConn.Close()

	log.Printf("Successfully connected to auth service")

	log.Printf("AUTH_SERVICE_ADDR: %s", os.Getenv("AUTH_SERVICE_ADDR"))

	authClient := pb.NewAuthServiceClient(authConn)

	// Initialize components
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)
	auth := authMiddleware.NewAuthMiddleware(authClient)

	// Setup router
	r := chi.NewRouter()

	//For metrics
	r.Handle("/metrics", promhttp.Handler())

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(authMiddleware.MetricsMiddleware)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.Authenticate)

		r.Post("/products", productHandler.Create)
		r.Get("/products/{id}", productHandler.Get)
		r.Get("/products", productHandler.List)
		r.Put("/products/{id}", productHandler.Update)
		r.Delete("/products/{id}", productHandler.Delete)
		r.Get("/user/products", productHandler.ListUserProducts)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting product service on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
