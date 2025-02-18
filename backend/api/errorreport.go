package api

import (
	"kochbuch-v2-backend/services"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ErrorReport struct {
	Url      string `json:"url" binding:"required"`
	Error    string `json:"error" binding:"required"`
	Severity string `json:"severity" binding:"required"`
}

// PostErrorReport handles the /errorreport endpoint
func PostErrorReport(c *gin.Context) {
	var report ErrorReport

	err := c.BindJSON(&report)
	if err != nil {
		log.Printf("Failed to parse input: %v", err)
		c.String(http.StatusBadRequest, "")
		return
	}

	query := "INSERT INTO `apilog`(`severity`, `reporter`, `host`, `agent`, `request_uri`, `message`) VALUES(?, 'Client', ?, ?, ?, ?)"
	message := report.Error
	host, _ := os.Hostname()
	agent := c.Request.Header.Get("user-agent")

	stmt, err := services.Db.Prepare(query)
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		c.String(http.StatusOK, "")
		return
	}

	_, err = stmt.Exec(report.Severity, host, agent, report.Url, message)
	if err != nil {
		log.Printf("Failed to execute statement: %v", err)
	}

	c.String(http.StatusOK, "")

}
