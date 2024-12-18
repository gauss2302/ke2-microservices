package main

import (
	"github.com/go-chi/cors"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gauss2302/testcommm/auth/internal/handler"
	"github.com/gauss2302/testcommm/auth/internal/pkg/jwt"
	"github.com/gauss2302/testcommm/auth/internal/service"
	pb_auth "github.com/gauss2302/testcommm/auth/proto/auth"
	pb_user "github.com/gauss2302/testcommm/auth/proto/user"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	// gRPC connection to User Service
	userConn, err := grpc.Dial(
		os.Getenv("USER_SERVICE_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer userConn.Close()

	// Create gRPC client
	userClient := pb_user.NewUserServiceClient(userConn)

	// Initialize services
	jwtMaker := jwt.NewJWTMaker(os.Getenv("JWT_PRIVATE_KEY"))
	authService := service.NewAuthService(rdb, userClient, jwtMaker)
	authHandler := handler.NewAuthHandler(authService)

	// gRPC server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb_auth.RegisterAuthServiceServer(grpcServer, authService)

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("starting gRPC server on :50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// HTTP server setup
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5555", "http://localhost:3000"}, // добавьте все нужные origin'ы
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // разрешаем все заголовки
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Get("/session", authHandler.GetSession)

	log.Printf("starting HTTP server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
