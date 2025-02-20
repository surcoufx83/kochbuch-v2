package main

import (
	"context"
	"kochbuch-v2-backend/api"
	"kochbuch-v2-backend/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Connect to MySQL database
	services.DbConnect()

	// Check connection to Nextcloud API
	services.NcConnect()

	// Load entities into cache
	services.LoadCategories(services.Db)
	services.LoadUnits(services.Db)
	services.LoadPublicRecipes(services.Db)

	// Set up Gin router
	router := gin.Default()
	router.Use(api.CheckValidHostnames())

	// Set up routes
	router.GET("/", api.GetIndex)
	router.GET("/categories", api.GetCategories)
	router.GET("/params", api.GetApiParams)
	router.POST("/errorreport", api.PostErrorReport)
	router.GET("/recipes", api.GetRecipes)
	router.GET("/units", api.GetUnits)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Handle signals to close database connection gracefully
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		log.Println("Server running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	sig := <-signalChan
	log.Printf("Received signal: %v. Shutting down server...", sig)

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped gracefully")
	services.Db.Close()
	log.Println("Database connection closed")

}
