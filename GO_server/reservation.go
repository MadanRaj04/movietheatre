package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterReservationRoutes(rg *gin.RouterGroup) {
	s := rg.Group("")
	{
		s.POST("/reservations", AuthMiddleware(), createReservation)
		s.GET("/reservations", AuthMiddleware(), listUserReservations)
		s.GET("/reservations/:id", AuthMiddleware(), getReservation)
	}
}

func createReservation(c *gin.Context) {
	var body struct {
		ShowtimeID uint  `json:"showtimeId" binding:"required"`
		Seats      []int `json:"seats" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := c.GetUint("userId")
	// fetch showtime
	var s Showtime
	if err := DB.First(&s, body.ShowtimeID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "showtime not found"})
		return
	}
	// basic seat collision check (simple): verify existing reservations for same showtime
	var existing []Reservation
	DB.Where("showtime_id = ?", body.ShowtimeID).Find(&existing)
	taken := map[int]bool{}
	for _, r := range existing {
		parts := strings.Split(r.Seats, ",")
		for _, p := range parts {
			if p == "" {
				continue
			}
			n, _ := strconv.Atoi(p)
			taken[n] = true
		}
	}
	for _, seat := range body.Seats {
		if taken[seat] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "seat already taken", "seat": seat})
			return
		}
	}
	// compute price
	total := float64(len(body.Seats)) * s.Price
	seatsStr := ""
	for i, v := range body.Seats {
		if i > 0 {
			seatsStr += ","
		}
		seatsStr += strconv.Itoa(v)
	}
	res := Reservation{UserID: userId, CinemaID: s.CinemaID, ShowtimeID: s.ID, Seats: seatsStr, TotalPrice: total}
	if err := DB.Create(&res).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// generate QR and send email
	qrContent := fmt.Sprintf("reservation:%d", res.ID)
	qrPath, _ := GenerateQRCodePNG(qrContent, fmt.Sprintf("res_%d.png", res.ID))
	res.QRCodePath = qrPath
	DB.Save(&res)
	// send mail (best effort, ignore errors)
	var user User
	DB.First(&user, userId)
	subject := "Your reservation"
	bodyHTML := fmt.Sprintf("<p>Reservation #%d created. Seats: %s</p>", res.ID, res.Seats)
	SendMail([]string{user.Email}, subject, bodyHTML)
	c.JSON(http.StatusCreated, res)
}

func listUserReservations(c *gin.Context) {
	userId := c.GetUint("userId")
	var r []Reservation
	DB.Where("user_id = ?", userId).Find(&r)
	c.JSON(http.StatusOK, r)
}

func getReservation(c *gin.Context) {
	id := c.Param("id")
	var r Reservation
	if err := DB.First(&r, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, r)
}
