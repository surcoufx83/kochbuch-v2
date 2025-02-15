package api

import (
	"kochbuch-v2-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUnits handles the /units endpoint
func GetUnits(c *gin.Context) {

	// Get units from cache
	units, etag := services.GetUnits()

	// Get Etag from request
	requestEtag := c.Request.Header.Get("If-None-Match")

	if requestEtag != "" && requestEtag == etag {
		// Etag matches, return 304 Not Modified
		c.Status(http.StatusNotModified)
	} else {
		// Set Etag header
		c.Header("Etag", etag)

		c.JSON(http.StatusOK, gin.H{
			"units": units,
			"etag":  etag,
		})
	}
}
