package main

import (
	"auth-service/internal/delivery/grpc"
	"auth-service/internal/delivery/http"
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/internal/usecase"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// init DB
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=123456 dbname=microservice_demo_auth port=5433 sslmode=disable TimeZone=Asia/Jakarta"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = db.AutoMigrate(
		&domain.User{},
	); err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}

	// init redis client
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
		log.Println("Continuing without Redis - token blacklisting will not work")
	}

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
	grpcHandler := grpc.NewGRPCHandler(authUseCase)

	// init gin router
	router := gin.Default()

	// register routes
	authHandler.RegisterRoutes(router)

	// channel get signal shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 2)

	// start HTTP server
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	go func() {
		log.Printf("HTTP server is running on :%s", httpPort)
		if err := router.Run(":" + httpPort); err != nil {
			errChan <- fmt.Errorf("Failed to start HTTP server: %v", err)
		}
	}()

	// start gRPC server
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	go func() {
		log.Printf("gRPC server is running on :%s", grpcPort)
		if err := grpcHandler.Serve(":" + grpcPort); err != nil {
			errChan <- fmt.Errorf("Failed to start gRPC server: %v", err)
		}
	}()

	// waiting signal
	select {
	case err := <-errChan:
		log.Printf("Error occured: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
	}

	// cleanup dan graceful shutdown can be add in here
	log.Println("Shutting down server...")
}
