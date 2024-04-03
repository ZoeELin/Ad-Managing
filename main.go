package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"

	"ad-proj/controllers"
	"ad-proj/database"
	"ad-proj/models"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type AdItem struct {
	Title string `json:"title"`
	EndAt string `json:"endAt"`
}

type AdItemslice struct {
	Items []AdItem
}

type AdResponse struct {
	Items   string `json:"item"`
	AdItems string `json:"adItems"`
}

var ads []models.Ad

func listAds(c *gin.Context) {
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

	// This is a slice that is compared to specified age, gender, countries, platforms
	var filteredAds []models.Ad
	for _, ad := range ads {
		if filterAd(ad, age, gender, country, platform) {
			filteredAds = append(filteredAds, ad)
		}
	}

	// Sort filtered ads by endAt (ASC)
	sort.SliceStable(filteredAds, func(i, j int) bool {
		return filteredAds[i].EndAt < filteredAds[j].EndAt
	})

	// Limit the number of ads to return
	if len(filteredAds) > limit {
		startIndex := offset
		endIndex := offset + limit
		if endIndex > len(filteredAds) {
			endIndex = len(filteredAds)
		}
		filteredAds = filteredAds[startIndex:endIndex]
	}

	fmt.Println("filteredAds", filteredAds)

	// Marshal the filtered ads into JSON and send the response
	var items []AdItem
	for _, ad := range filteredAds {
		item := AdItem{
			Title: ad.Title,
			EndAt: ad.EndAt,
		}
		items = append(items, item)
	}
	response := map[string]interface{}{"items": items}

	c.JSON(http.StatusOK, response)
}

func filterAd(ad models.Ad, age int, gender string, country string, platform string) bool {
	for _, condition := range ad.Conditions {
		if (age == 0 || age >= condition.AgeStart) &&
			(age == 0 || age <= condition.AgeEnd) &&
			(gender == "" || gender == condition.Gender) &&
			(country == "" || isInSlice(country, condition.Country)) &&
			(platform == "" || isInSlice(platform, condition.Platform)) {
			return true
		}
	}
	return false
}

// Check if a slice contains a string that query
func isInSlice(text string, sliceText []string) bool {
	for _, c := range sliceText {
		if c == text {
			return true
		}
	}
	return false
}

func main() {
	// Get parameter from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}

	// Set the parameters to dsn
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	user := "postgres"
	dbname := "postgres"
	pass := os.Getenv("PASSWORD")
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)

	// Controller to initiallize the database
	initFlag := flag.Bool("init", false, "Initialize the database")
	flag.Parse()

	// if add -init, initiallize db
	if *initFlag {
		database.DbInit(dsn)
		database.DatasetInit(dsn)
		fmt.Println("Initialization completed.")
	}

	r := gin.Default()

	database.ConnectDatabase(dsn)

	r.POST("/api/v1/ad", controllers.CreateAd)
	r.GET("/api/v1/ad", listAds)

	fmt.Println("Server is running on port 8088")
	r.Run(":8088")
}
