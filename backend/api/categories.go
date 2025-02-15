package api

import (
	"kochbuch-v2-backend/cache"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCategories handles the /categories endpoint
func GetCategories(c *gin.Context) {

	// Get categories from cache
	categories, etag := cache.GetCategories()

	// Get Etag from request
	requestEtag := c.Request.Header.Get("If-None-Match")

	if requestEtag == etag {
		// Etag matches, return 304 Not Modified
		c.Status(http.StatusNotModified)
	} else {
		// Set Etag header
		c.Header("Etag", etag)

		c.JSON(http.StatusOK, gin.H{
			"categories": categories,
		})
	}
}
