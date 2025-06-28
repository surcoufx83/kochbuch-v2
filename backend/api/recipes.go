package api

import (
	"fmt"
	"kochbuch-v2-backend/services"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

/* func GetRecipes(c *gin.Context) {

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
} */

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
		log.Printf("  a> Recipe not found: %v", err)
		c.Status(http.StatusNotFound)
		return
	}

	_, picture, err := services.GetPicture(recipe, uint32(pictureid))
	if err != nil || picture.Id == 0 || picture.FileName != c.Param("filename") {
		log.Printf("  b> Picture not found: %v", err)
		c.Status(http.StatusNotFound)
		return
	}

	c.File(picture.FullPath)

}

func GetRecipeThbPicture(c *gin.Context) {
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
		log.Printf("  a> Recipe not found: %v", err)
		c.Status(http.StatusNotFound)
		return
	}

	_, picture, err := services.GetPicture(recipe, uint32(pictureid))
	if err != nil || picture.Id == 0 {
		log.Printf("  b> Picture not found: %v", err)
		c.Status(http.StatusNotFound)
		return
	}

	thbsize, err := strconv.Atoi(c.Param("thbsize"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	for _, size := range picture.Dimension.GeneratedSizes {
		if size == thbsize {
			folder := filepath.Dir(picture.FullPath)
			basename, ext := services.GetBasenameAndExtension(picture.FileName)
			resizedFilename := filepath.Join(folder, fmt.Sprintf("%s_%d%s", basename, size, ext))
			c.File(resizedFilename)
			return
		}
	}

	log.Printf("  c> Size not found: %v", thbsize)
	c.Status(http.StatusNotFound)

}
