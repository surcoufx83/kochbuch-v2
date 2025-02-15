package api

import (
	"kochbuch-v2-backend/cache"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCategories handles the /categories endpoint
func GetCategories(c *gin.Context) {

	// Get categories from cache
	categories, etag := cache.GetCategories()

	requestEtag := c.Request.Header.Get("If-None-Match")

	log.Println("Request Etag: ", requestEtag)
	log.Println("        Etag: ", etag)

	if requestEtag == etag {
		c.Status(http.StatusNotModified)
	} else {
		// Set Etag header
		c.Header("Etag", etag)

		c.JSON(http.StatusOK, gin.H{
			"categories": categories,
		})
	}
}
