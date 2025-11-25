package main

import (
	"log"
	"net"

	"github.com/Kars07/rest-grpc-demo/database"
	grpcServer "github.com/Kars07/rest-grpc-demo/grpc/server"
	"github.com/Kars07/rest-grpc-demo/handlers"
	pb "github.com/Kars07/rest-grpc-demo/proto"
	"github.com/Kars07/rest-grpc-demo/repository"
	"github.com/Kars07/rest-grpc-demo/service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize database
	database.ConnectDatabase()
	db := database.GetDB()

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// Start REST API server in a goroutine
	go startRESTServer(userService)

	// Start gRPC server in the main goroutine
	startGRPCServer(userService)
}

func startRESTServer(userService *service.UserService) {
	// Create Gin router
	router := gin.Default()

	// Initialize handler
	userHandler := handlers.NewUserHandler(userService)
	userHandler.RegisterRoutes(router)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Starting REST API server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start REST server: %v", err)
	}
}

func startGRPCServer(userService *service.UserService) {
	// Create listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	grpcSrv := grpc.NewServer()

	// Register service
	userServer := grpcServer.NewUserServer(userService)
	pb.RegisterUserServiceServer(grpcSrv, userServer)

	// Enable reflection for grpcurl
	reflection.Register(grpcSrv)

	log.Println("Starting gRPC server on :50051...")
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
