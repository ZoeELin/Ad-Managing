package database

import (
	"ad-proj/models"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

var Db *gorm.DB

func ConnectDatabase(dsn string) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Successfully connected to database and created a ads table!")
	Db = db
}

// Initialize the database and create a table
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

	fmt.Println("Database initialization completed successfully.")
}

// Insert 100 random data
func DatasetInit(dsn string) {
	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Insert random data
	err = InsertRandomData(db, 100) // Insert 100 random records
	if err != nil {
		log.Fatalf("Failed to insert random data: %v", err)
	}
}

// InsertRandomData inserts random data into the AdsColumn table
func InsertRandomData(db *gorm.DB, count int) error {
	for i := 0; i < count; i++ {

		ageStart, ageEnd := generateRandomAge()

		ad := models.AdsColumn{
			Title:    fmt.Sprintf("AD-%02d", i+1),
			StartAt:  generateRandomTime(),
			EndAt:    generateRandomTime(),
			AgeStart: ageStart,
			AgeEnd:   ageEnd,
			Gender:   randomGender(),
			Country:  randomCountries(),
			Platform: randomPlatforms(),
		}
		if err := db.Create(&ad).Error; err != nil {
			return err
		}
	}
	return nil
}

// InsertData data into the AdsColumn table
func InsertData(ad models.Ad) error {
	for _, condition := range ad.Conditions {
		startAtTime, err := time.Parse("2006-01-02T15:04:05.000Z", ad.StartAt)
		if err != nil {
			return err
		}

		endAtTime, err := time.Parse("2006-01-02T15:04:05.000Z", ad.EndAt)
		if err != nil {
			return err
		}

		ad := models.AdsColumn{
			Title:    ad.Title,
			StartAt:  startAtTime,
			EndAt:    endAtTime,
			AgeStart: condition.AgeStart,
			AgeEnd:   condition.AgeEnd,
			Gender:   condition.Gender,
			Country:  strings.Join(condition.Country, ","),
			Platform: strings.Join(condition.Platform, ","),
		}

		if err := Db.Create(&ad).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetActiveAdsCount returns the count of active ads based on the current time
func GetActiveAdsCount() (int64, error) {
	var count int64

	// query := Db.Model(&models.AdsColumn{})

	// Run a query to calculate the number of active ads here
	err := Db.Model(&models.AdsColumn{}).Where("start_at < ? AND end_at > ?", time.Now(), time.Now()).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// SelectData retrieves ads from the dat[error] unsupported data type: &[]abase based on specified criteria
func SelectData(filteredAds *[]models.AdsColumn, offset int, limit int, age int, gender, country, platform string) error {
	query := Db.Model(&models.AdsColumn{})

	// Add conditions based on the specified criteria
	if age != 0 {
		query = query.Where("age_start <= ? AND age_end >= ?", age, age)
	}
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if country != "" {
		query = query.Where("country LIKE ?", "%"+country+"%")
	}
	if platform != "" {
		query = query.Where("platform LIKE ?", "%"+platform+"%")
	}

	// Convert EndAt to string for sorting
	query = query.Order("end_at ASC").Offset(offset).Limit(limit)

	// Execute the query and scan the results into the provided slice
	if err := query.Find(&filteredAds).Error; err != nil {
		return err
	}

	return nil
}

// Generate the age, notice that AgeStart smaller or equal to AgeEnd
func generateRandomAge() (int, int) {
	ageStart := rand.Intn(100) + 1
	ageEnd := rand.Intn(100) + 1
	if ageStart >= ageEnd {
		// Swap AgeStart and AgeEnd if AgeStart is greater than or equal to AgeEnd
		ageStart, ageEnd = ageEnd, ageStart
	}
	return ageStart, ageEnd
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
