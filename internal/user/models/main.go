package models

import "time"

// User represents a user in the system
type User struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	// Stocks      []string           `json:"stocks" db:"stocks"`
	// AlertLevels map[string]float64 `json:"alert_levels" db:"alert_levels"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
