package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"product-service/internal/delivery/grpc"
	"product-service/internal/delivery/http"
	"product-service/internal/domain"
	"product-service/internal/repository"
	"product-service/internal/usecase"
	"syscall"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// init DB
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=123456 dbname=microservice_demo_product port=5433 sslmode=disable TimeZone=Asia/Jakarta"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = db.AutoMigrate(&domain.Product{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// init repository
	productRepo := repository.NewProductRepository(db)

	// init usecase
	productUseCase := usecase.NewProductUseCase(productRepo)

	// init HTTP handler
	productHandler := http.NewProductHandler(productUseCase)

	// init gRPC handler
	grpcHandler := grpc.NewGRPCProductHandler(productUseCase)

	// init gin router
	router := gin.Default()
	router.Use(CorsMiddleware())

	// register routes
	productHandler.RegisterRoutes(router)

	// channel signal shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// chan err server
	errChan := make(chan error, 2)

	// start HTTP server
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8081"
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", httpPort)
		if err := router.Run(":" + httpPort); err != nil {
			errChan <- fmt.Errorf("Failed to start HTTP server: %v", err)
		}
	}()

	// start gRPC server
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}

	go func() {
		log.Printf("Starting gRPC server on port %s", grpcPort)
		if err := grpcHandler.Serve(":" + grpcPort); err != nil {
			errChan <- fmt.Errorf("Failed to start gRPC server: %v", err)
		}
	}()

	// wait for shutdown signal or error
	select {
	case err := <-errChan:
		log.Printf("Error occured: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
	}

	log.Println("Shutting down server...")

	// cleanup code (e.g., close database connection)
	if db, err := db.DB(); err == nil {
		db.Close()
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
