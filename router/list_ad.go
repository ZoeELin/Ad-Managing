// router/list_ad.go

package router

import (
	"ad-proj/database"
	"ad-proj/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HTTP GET request, api that can get ads from database
func ListAds(c *gin.Context) {
	// Parse query parameters
	offsetStr := c.DefaultQuery("offset", "1")
	limitStr := c.DefaultQuery("limit", "5")
	ageStr := c.Query("age")
	gender := c.Query("gender")
	country := c.Query("country")
	platform := c.Query("platform")

	// Convert strings to integers
	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)
	age, _ := strconv.Atoi(ageStr)

	// This is a slice that is compared to specified age, gender, country, platform
	var filteredAds []models.AdsColumn

	// SelectData retrieves ads from the database
	err := database.SelectData(&filteredAds, offset, limit, age, gender, country, platform)
	if err != nil {
		panic("Failed to retrieve ads from database")
	}

	// Reorganize data from db to response
	var items []models.AdItem
	for _, ad := range filteredAds {

		formattedTime := ad.EndAt.Format("2006-01-02T15:04:05.000Z")

		item := models.AdItem{
			Title: ad.Title,
			EndAt: formattedTime,
		}
		items = append(items, item)
	}
	response := map[string]interface{}{"items": items}

	c.JSON(http.StatusOK, response)
}
