package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Ad represents the structure of an advertisementgit
type Ad struct {
	Title      string      `json:"title"`
	StartAt    string      `json:"startAt"`
	EndAt      string      `json:"endAt"`
	Conditions []Condition `json:"conditions,omitempty"`
}

// Condition represents a condition for targeting advertisements
type Condition struct {
	AgeStart  int      `json:"ageStart,omitempty"`
	AgeEnd    int      `json:"ageEnd,omitempty"`
	Gender    string   `json:"gender,omitempty"`
	Countries []string `json:"countries,omitempty"`
	Platforms []string `json:"platforms,omitempty"`
}

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

var ads []Ad

// HTTP POST request, input 2 parameters
func createAd(c *gin.Context) {
	// A variable to store new ad which is encoded
	var newAd Ad

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
	var filteredAds []Ad
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

func filterAd(ad Ad, age int, gender string, country string, platform string) bool {
	// fmt.Printf("condition.Countries, countries", condition.Countries, country)
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

// Generate a random time string in the format "2006-01-02T15:04:05.000Z"
func generateRandomTime() string {
	min := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC).Unix()
	randomUnix := rand.Int63n(max-min) + min
	randomTime := time.Unix(randomUnix, 0)
	return randomTime.Format("2006-01-02T15:04:05.000Z")
}

// Generate a random gender string "M" or "F"
func randomGender() string {
	genders := []string{"M", "F"}
	return genders[rand.Intn(len(genders))]
}

// Generate random countries
func randomCountries() []string {
	countries := []string{"TW", "JP", "US", "CA", "UK", "AU"}
	numCountries := rand.Intn(len(countries))
	rand.Shuffle(len(countries), func(i, j int) {
		countries[i], countries[j] = countries[j], countries[i]
	})
	return countries[:numCountries]
}

// Generate random platforms
func randomPlatforms() []string {
	platforms := []string{"android", "ios", "web"}
	numPlatforms := rand.Intn(len(platforms))
	rand.Shuffle(len(platforms), func(i, j int) {
		platforms[i], platforms[j] = platforms[j], platforms[i]
	})
	return platforms[:numPlatforms]
}

func main() {

	// Seed the random number generator for consistent results
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		ad := Ad{
			Title:   fmt.Sprintf("AD %02d", i+1),
			StartAt: generateRandomTime(),
			EndAt:   generateRandomTime(),
			Conditions: []Condition{
				{
					AgeStart:  rand.Intn(100) + 1,
					AgeEnd:    rand.Intn(100) + 1,
					Gender:    randomGender(),
					Countries: randomCountries(),
					Platforms: randomPlatforms(),
				},
			},
		}
		ads = append(ads, ad)
	}

	// Convert ads slice to JSON
	adsJSON, err := json.MarshalIndent(ads, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Write JSON data to a file
	file, err := os.Create("ads_data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(adsJSON)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ads data generated and saved to ads_data.json")

	r := gin.Default()

	r.POST("/api/v1/ad", createAd)
	r.GET("/api/v1/ad", listAds)

	fmt.Println("Server is running on port 8088")
	r.Run(":8088")
}
