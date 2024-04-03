// models/ad.go

package models

import (
	"time"
)

// Ad represents the structure of an advertisement
type Ad struct {
	Title      string      `json:"title"`
	StartAt    string      `json:"startAt"`
	EndAt      string      `json:"endAt"`
	Conditions []Condition `json:"conditions,omitempty"`
}

// Condition represents the conditions for displaying an advertisement
type Condition struct {
	AgeStart  int      `json:"ageStart"`
	AgeEnd    int      `json:"ageEnd"`
	Gender    string   `json:"gender"`
	Countries []string `json:"countries"`
	Platforms []string `json:"platforms"`
}

type AdsColumn struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"column:title"`
	StartAt   time.Time `gorm:"column:start_at"`
	EndAt     time.Time `gorm:"column:end_at"`
	AgeStart  int       `gorm:"column:age_start"`
	AgeEnd    int       `gorm:"column:age_end"`
	Gender    string    `gorm:"column:gender"`
	Countries string    `gorm:"column:countries"`
	Platforms string    `gorm:"column:platforms"`
}