// user-service/cmd/main.go
package main

import (
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gauss2302/testcommm/user/internal/domain/entity"
	"github.com/gauss2302/testcommm/user/internal/repository"
	"github.com/gauss2302/testcommm/user/internal/service"
	pb "github.com/gauss2302/testcommm/user/proto"
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
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
	}
	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	log.Printf("Successfully connected to database")

	// Auto migrate
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// gRPC server setup
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userService)

	log.Printf("Starting gRPC server on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
