// controllers/ads.go

package controllers

import (
	"ad-proj/database"
	"ad-proj/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTP POST request, api that can insert ad info to database
func CreateAd(c *gin.Context) {
	// A variable to store new ad which is encoded
	var newAd models.Ad

	// Decode the newAd, return err if failed
	if err := c.ShouldBindJSON(&newAd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the newAd is in correct JSON formation
	if newAd.Title == "" || newAd.StartAt == "" || newAd.EndAt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect formation. Title, StartAt, and EndAt are required fields."})
		return
	}

	for _, condition := range newAd.Conditions {

		// Check if AgeStart and AgeEnd are provided, and if provided, check their range
		if condition.AgeStart != 0 && (condition.AgeStart < 1 || condition.AgeStart > 100) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect formation. AgeStart needs to be in the range from 1 to 100."})
			return
		}
		if condition.AgeEnd != 0 && (condition.AgeEnd < 1 || condition.AgeEnd > 100) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect formation. AgeEnd needs to be in the range from 1 to 100."})
			return
		}

		// Check Gender
		if condition.Gender != "" && condition.Gender != "F" && condition.Gender != "M" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect formation. Gender needs to be M or F"})
			return
		}
	}

	database.InsertData(newAd)
	c.JSON(http.StatusCreated, gin.H{"message": "Post the ad successfully.", "data": newAd})
}
