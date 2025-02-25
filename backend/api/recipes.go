package api

import (
	"kochbuch-v2-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRecipes(c *gin.Context) {

	etag := services.GetRecipesEtag()

	// Get Etag from request
	requestEtag := c.Request.Header.Get("If-None-Match")

	if requestEtag == etag {
		// Etag matches, return 304 Not Modified
		c.Status(http.StatusNotModified)
	} else {
		// Get categories from cache
		recipes, _ := services.GetRecipes(c)

		// Set Etag header
		c.Header("Etag", etag)

		c.JSON(http.StatusOK, gin.H{
			"recipes": recipes,
		})
	}
}

func GetRecipePicture(c *gin.Context) {
	projectid, err := strconv.Atoi(c.Param("projectid"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	pictureid, err := strconv.Atoi(c.Param("pictureid"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	recipe, err := services.GetRecipe(uint32(projectid), c)
	if err != nil || recipe.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	for _, pic := range recipe.Pictures {
		if pic.Id == uint32(pictureid) && pic.FileName == c.Param("filename") {
			c.File(pic.FullPath)
			return
		}
	}

	c.Status(http.StatusNotFound)

}
