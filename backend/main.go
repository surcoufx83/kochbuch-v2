package main

import (
	"context"
	"fmt"
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
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins, adjust this for production
	},
}

func main() {

	// Connect to MySQL database
	services.DbConnect()

	// Check connection to Nextcloud API
	services.NcConnect()

	services.Locales = []string{"de", "en", "fr"}
	services.ThumbnailSizes = []int{300, 400, 500, 600, 800}

	// Check connection to AI assistant
	go services.AiConnect()

	// Load entities into cache
	services.LoadCategories()
	services.LoadUnits()
	services.LoadRecipes()

	// Run Image Thumbnail generator
	go services.ThbAutoGenerator()

	// Set up Gin router
	router := gin.Default()
	router.Use(api.CheckValidHostnames())

	// Set up routes
	router.GET("/", api.GetIndex)
	// router.GET("/categories", api.GetCategories)
	router.POST("/errorreport", api.PostErrorReport)
	/* 	router.POST("/login", api.PostOauth2Login)
	   	router.POST("/logout", api.PostLogout)
	   	router.GET("/me", api.GetMyProfile) */
	// router.GET("/params", api.GetAppParams)
	// router.GET("/recipes", api.GetRecipes)
	// router.GET("/units", api.GetUnits)
	router.GET("/ws", services.OnWebsocketConnect)

	media := router.Group("/media", CacheMiddleware(2592000))
	{
		media.GET("/uploads/:projectid/:pictureid/:filename", api.GetRecipePicture)
		media.GET("/uploads/:projectid/:pictureid/thb/:thbsize/:filename", api.GetRecipeThbPicture)
	}

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
	services.DbCloseStmts()
	services.Db.Close()
	log.Println("Database connection closed")

}

func CacheMiddleware(maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		c.Next()
	}
}
