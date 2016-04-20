package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	defaultPort = "8080"
)

type trackingRequest struct {
	vendor string
	trackingNumber string
}

func main() {
	port := getPort()
	engine := gin.Default()

	engine.GET("/uni/:tracking_number", func(c *gin.Context) {
		trackingNumber := c.Param("tracking_number")
		c.JSON(http.StatusOK, gin.H{
			"tracking_number": trackingNumber,
		})
	})

	engine.Run(":" + port)
}

func getPort() string {
	textPort := os.Getenv("PORT")
	_, err := strconv.Atoi(textPort)
	if err != nil {
		return defaultPort
	}
	return textPort
}
