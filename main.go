package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort = "8080"
)

func main() {
	port := getPort()
	engine := gin.Default()

	engine.GET("/uni/:tracking_number", func(c *gin.Context) {
		trackingNumber := c.Param("tracking_number")
		summary, err := getTrackingSummary(trackingNumber)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"tracking_summary": summary,
			})
		}
	})

	engine.Run(":" + port)
}

func getTrackingSummary(trackingNumber string) (string, error) {
	// Get from web.
	getRequest := "http://production.shippingapis.com/ShippingAPI.dll?API=TrackV2&" +
		"XML=%3CTrackRequest%20USERID%3D%22315INDIV8018%22%3E%3CTrackID%20ID%3D%22" +
		trackingNumber +
		"%22%3E%3C%2FTrackID%3E%3C/TrackRequest%3E"
	getResponse, err := http.Get(getRequest)
	if err != nil {
		return "", err
	}
	defer getResponse.Body.Close()
	body, err := ioutil.ReadAll(getResponse.Body)
	if err != nil {
		return "", err
	}

	// Parse xml.
	xmlString := string(body)
	reader := strings.NewReader(xmlString)
	decoder := xml.NewDecoder(reader)
	summaryStarted := false
	for token, err := decoder.Token(); err == nil; token, err = decoder.Token() {
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "TrackSummary" {
				summaryStarted = true
			}
		case xml.CharData:
			if summaryStarted {
				return string([]byte(t)), nil
			}
		}
	}

	return string(body), nil
}

func getPort() string {
	textPort := os.Getenv("PORT")
	_, err := strconv.Atoi(textPort)
	if err != nil {
		return defaultPort
	}
	return textPort
}
