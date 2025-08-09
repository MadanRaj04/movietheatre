package main

import (
	"time"

	"gorm.io/gorm"
)

// User represents an app user
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `json:"name" gorm:"size:200"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:200"`
	Password  string         `json:"-"`
	Username  string         `json:"username" gorm:"size:50"`
	Role      string         `json:"role" gorm:"size:50"`
}

// Movie represents a film
type Movie struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Title       string    `json:"title" gorm:"size:300"`
	Description string    `json:"description" gorm:"type:text"`
	Duration    int       `json:"duration"`
	Poster      string    `json:"poster"`
}

// Cinema represents a cinema/hall
type Cinema struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name" gorm:"size:200"`
	Location  string    `json:"location"`
	Seats     int       `json:"seats"` // number of seats
}

// Showtime for a movie in a cinema
type Showtime struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	MovieID   uint      `json:"movieId"`
	CinemaID  uint      `json:"cinemaId"`
	Start     time.Time `json:"start"`
	Price     float64   `json:"price"`
}

// Reservation stores booked seats
type Reservation struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	UserID     uint      `json:"userId"`
	CinemaID   uint      `json:"cinemaId"`
	ShowtimeID uint      `json:"showtimeId"`
	Seats      string    `json:"seats"` // comma-separated seat numbers
	TotalPrice float64   `json:"totalPrice"`
	QRCodePath string    `json:"qrcodePath"`
}
