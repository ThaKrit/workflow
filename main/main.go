package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-api/items"
	"go-api/middle"
	"go-api/users"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Define timeout and delay durations
	const shutdownDelay = 3 * time.Second
	const gracefulShutdownTimeout = 3 * time.Second

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get server port from environment variables
	port := os.Getenv("PORT")

	// Setup database connection
	dsn := os.Getenv("POSTGRES_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// // Migrate the schema (optional, but useful to keep database up to date)
	// if err := db.AutoMigrate(&model.Item{}, &model.User{}); err != nil {
	// 	log.Fatal("Failed to migrate database schema:", err)
	// }

	// Initialize repository, service, and controller for items
	itemRepo := items.NewRepository(db)
	itemService := items.NewService(itemRepo)
	itemController := items.NewController(itemService)

	// Initialize repository, service, and controller for users
	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)
	userController := users.NewUserController(userService)

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // อนุญาตทุก origin
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// LOGIN routes
	r.POST("/users/login", userController.UserLogin)
	r.POST("/users", userController.CreateUser)

	// ดัก Verify JWT
	r.Use(middle.Guard(os.Getenv("JWT_SECRET")))

	// Register user routes
	r.GET("/users", userController.GetUsers)
	r.GET("/users/:id", userController.GetUser)
	r.PUT("/users/:id", userController.UpdateUser)
	r.PATCH("/users/:id", userController.PatchUser)
	r.DELETE("/users/:id", userController.DeleteUser)

	// Register item routes
	r.POST("/items", itemController.CreateItem)
	r.GET("/items", itemController.GetItems)
	r.GET("/items/:id", itemController.GetItem)
	r.PUT("/items/:id", itemController.UpdateItem)
	r.PATCH("/items/:id", itemController.PatchItem)
	r.DELETE("/items/:id", itemController.DeleteItem)

	// Create an HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		fmt.Printf("Server running on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen: %v\n", err)
		}
	}()

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Wait until a signal is received

	// Log message with gracefulShutdownTimeout
	log.Printf("Received shutdown signal. Delaying shutdown for %v...", gracefulShutdownTimeout)

	// Countdown before shutdown
	for i := int(shutdownDelay / time.Second); i > 0; i-- {
		fmt.Printf("%d...\n", i)
		time.Sleep(1 * time.Second)
	}

	log.Println("Shutting down server...")

	// Graceful shutdown with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
