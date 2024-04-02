package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"

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

// HTTP POST request, input 2 parameters
func createAd(c *gin.Context) {
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

		// Check Gender if validate it
		if condition.Gender != "" && condition.Gender != "F" && condition.Gender != "M" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect formation. Gender needs to be M or F"})
			return
		}
	}

	ads = append(ads, newAd)
	c.JSON(http.StatusCreated, gin.H{"message": "Post the ad successfully.", "ads": ads})
}

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
			(country == "" || isInSlice(country, condition.Countries)) &&
			(platform == "" || isInSlice(platform, condition.Platforms)) {
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
	err := godotenv.Load() //by default, it is .env so we don't have to write
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}
	//we read our .env file
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	user := "postgres"
	dbname := "postgres"
	pass := os.Getenv("PASSWORD")
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)

	// controller to initiallize the database
	initFlag := flag.Bool("init", false, "Initialize the database")
	flag.Parse()

	// if add -init, initiallize
	if *initFlag {
		database.DbInit(dsn)
		fmt.Println("Initialization completed.")
	} else {
		fmt.Println("No initialization performed.")
	}

	database.ConnectDatabase(dsn)

	r := gin.Default()

	r.POST("/api/v1/ad", createAd)
	r.GET("/api/v1/ad", listAds)

	fmt.Println("Server is running on port 8088")
	r.Run(":8088")
}
