package database

import (
	"ad-proj/models"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq" // don't forget to add it. It doesn't be added automatically
)

var Db *gorm.DB //created outside to make it global.

// make sure your function start with uppercase to call outside of the directory.
func ConnectDatabase(dsn string) {

	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Successfully connected to database and created a ads table!")
}

// DbInit initializes the database and inserts random data
func DbInit(dsn string) {
	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate tables
	err = db.AutoMigrate(&models.AdsColumn{})
	if err != nil {
		log.Fatalf("Failed to auto migrate tables: %v", err)
	}

	// Insert random data
	err = InsertRandomData(db, 100) // Insert 100 random records
	if err != nil {
		log.Fatalf("Failed to insert random data: %v", err)
	}

	fmt.Println("Database initialization completed successfully.")
}

// InsertRandomData inserts random data into the AdsColumn table
func InsertRandomData(db *gorm.DB, count int) error {
	for i := 0; i < count; i++ {
		ad := models.AdsColumn{
			Title:     fmt.Sprintf("AD %02d", i+1),
			StartAt:   generateRandomTime(),
			EndAt:     generateRandomTime(),
			AgeStart:  rand.Intn(100) + 1,
			AgeEnd:    rand.Intn(100) + 1,
			Gender:    randomGender(),
			Countries: randomCountries(),
			Platforms: randomPlatforms(),
		}
		if err := db.Create(&ad).Error; err != nil {
			return err
		}
	}
	return nil
}

// Generate a random time between now and 1 year from now
func generateRandomTime() time.Time {
	min := time.Now()
	max := min.AddDate(1, 0, 0)
	delta := max.Unix() - min.Unix()
	sec := rand.Int63n(delta) + min.Unix()
	return time.Unix(sec, 0)
}

// Generate a random gender string "M" or "F"
func randomGender() string {
	genders := []string{"M", "F"}
	return genders[rand.Intn(len(genders))]
}

// Generate random countries
func randomCountries() string {
	countries := []string{"TW", "JP", "US", "CA", "UK", "AU"}
	return countries[rand.Intn(len(countries))]
}

// Generate random platforms
func randomPlatforms() string {
	platforms := []string{"android", "ios", "web"}
	return platforms[rand.Intn(len(platforms))]
}
