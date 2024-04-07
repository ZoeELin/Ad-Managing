// controllers/ads.go

package controllers

import (
	"ad-proj/database"
	"ad-proj/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

const dailyAdsLimit = 3000
const activeAdsLimit = 1000

var dailyAdsCounter int

// HTTP POST request, api that can insert ad info to database
func CreateAd(c *gin.Context) {
	// Check if the maximum number of new ads per day has been exceeded
	if dailyAdsCounter >= dailyAdsLimit {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Daily ads creation limit reached. Please try again tomorrow."})
		return
	}

	// Check if the number of active ads exceeds the limit
	activeAdsCount, err := database.GetActiveAdsCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check active ads count."})
		return
	}
	if activeAdsCount >= activeAdsLimit {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Active ads limit reached. Please try again later."})
		return
	}

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

	// Add a counter for the number of ads added per day when an ad is successfully added.
	dailyAdsCounter++
	c.JSON(http.StatusCreated, gin.H{"message": "Post the ad successfully.", "data": newAd})
}
