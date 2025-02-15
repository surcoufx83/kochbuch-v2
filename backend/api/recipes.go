package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRecipes(c *gin.Context) {

	c.String(http.StatusNoContent, "")
}
