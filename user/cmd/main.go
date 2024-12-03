// user-service/cmd/main.go
package main

import (
	"github.com/gauss2302/testcommm/user/config"
	"github.com/gauss2302/testcommm/user/internal/domain/entity"
	"github.com/gauss2302/testcommm/user/internal/repository"
	"github.com/gauss2302/testcommm/user/internal/service"
	database "github.com/gauss2302/testcommm/user/pkg/databse"
	"github.com/gauss2302/testcommm/user/pkg/server"
	"log"
)

func main() {
	// Load Config
	cfg := config.Load()

	// Set up the database
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// Initialize and start gRPC server
	srv := server.NewServer(userService)
	log.Printf("Starting gRPC server on :%s", cfg.GRPC.Port)
	if err := srv.Start(cfg.GRPC); err != nil {
		log.Fatal(err)
	}
}
