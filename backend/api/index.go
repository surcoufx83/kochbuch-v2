package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIndex handles the root endpoint
func GetIndex(c *gin.Context) {
	c.String(http.StatusNoContent, "")
}
