package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetLoginParams handles the /api/login/params
func GetLoginParams(c *gin.Context) {

	c.String(http.StatusNoContent, "")
}
