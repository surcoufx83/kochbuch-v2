package main

import (
	"context"
	"kochbuch-v2-backend/api"
	"kochbuch-v2-backend/types"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB

	cacheMutex sync.RWMutex

	categoriesCache   map[uint16]types.Category
	categoriesEtag    time.Time
	categoriesEtagStr string
)

func main() {

	// Connect to MySQL database
	db_connect()

	// Load categories into cache at startup
	load_categories()

	// Set up Gin router
	router := gin.Default()
	router.GET("/", api.GetIndex)

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
	db.Close()
	log.Println("Database connection closed")

}

func db_connect() {
	dbuser := os.Getenv("DB_User")
	dbpassword := os.Getenv("DB_Password")
	dbhost := os.Getenv("DB_Host")
	dbport := os.Getenv("DB_Port")
	dbname := os.Getenv("DB_Name")
	tz := os.Getenv("TZ")

	if dbhost == "" {
		dbhost = "localhost"
	}
	if dbport == "" {
		dbport = "3306"
	}
	if tz == "" {
		tz = "UTC"
	}

	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?parseTime=true&loc=" + tz
	var err error

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Connected to database!")
}

func load_categories() {
	query := "SELECT * FROM categoryitemsview"
	var categories []types.Category

	err := db.Select(&categories, query)
	if err != nil {
		log.Fatalf("Failed to load categories: %v", err)
	}

	// Build cache
	cacheMutex.Lock()
	categoriesCache := make(map[uint16]types.Category)
	for _, category := range categories {
		categoriesCache[category.ItemId] = category
		if category.ItemModified.After(categoriesEtag) {
			categoriesEtag = category.ItemModified
		} else if category.CatModified.After(categoriesEtag) {
			categoriesEtag = category.CatModified
		}
	}
	categoriesEtagStr = categoriesEtag.Format(time.RFC3339)
	cacheMutex.Unlock()
	log.Printf("Loaded %d categories into cache", len(categories))
	log.Printf("Categories cache ETag: %v", categoriesEtagStr)
}
