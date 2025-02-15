package api

import (
	"kochbuch-v2-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRecipes(c *gin.Context) {

	// Get categories from cache
	recipes, etag := services.GetPublicRecipes()

	// Get Etag from request
	requestEtag := c.Request.Header.Get("If-None-Match")

	if requestEtag == etag {
		// Etag matches, return 304 Not Modified
		c.Status(http.StatusNotModified)
	} else {
		// Set Etag header
		c.Header("Etag", etag)

		c.JSON(http.StatusOK, gin.H{
			"recipes": recipes,
			"etag":    etag,
		})
	}
}
