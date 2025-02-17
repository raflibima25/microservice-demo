package main

import (
	"auth-service/internal/delivery/http"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/internal/usecase"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// init DB
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=123456 dbname=auth_service port=5433 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// init redis client
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// init repository
	userRepo := repository.NewUserRepository(db)

	// init services
	jwtSecretKey := os.Getenv("JWT_SECRET")
	if jwtSecretKey == "" {
		jwtSecretKey = "rahasia"
	}
	tokenService := service.NewJwtTokenService(jwtSecretKey, redisClient)

	// init use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenService)

	// init HTTP handler
	authHandler := http.NewAuthHandler(authUseCase)

	// init gRPC handler
	grpcHandler := grpc.NewAuthHandler(authUseCase)

	// init gin router
	router := gin.Default()

	// register routes
	authHandler.RegisterRoutes(router)

	// start HTTP server
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	go func() {
		if err := router.Run(":" + httpPort); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// start gRPC server
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	go func() {
		if err := grpcHandler.Serve(":" + grpcPort); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()
}
