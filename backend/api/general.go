package api

import (
	"kochbuch-v2-backend/services"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func CheckValidHostnames() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !services.ValidDomains[c.Request.Host] {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Next()
	}
}

// GetLoginParams handles the /api/params endpoint
func GetApiParams(c *gin.Context) {
	state, params, err := services.GetNcLoginParams(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	url, _ := url.Parse(c.Request.Host)
	if state != "" {
		c.SetCookie("session", state, 2592000, "/api", url.Hostname(), true, true)
	}

	c.JSON(http.StatusOK, params)
}

// GetIndex handles the root endpoint
func GetIndex(c *gin.Context) {
	c.String(http.StatusNoContent, "")
}
