package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterShowtimeRoutes(rg *gin.RouterGroup) {
	s := rg.Group("")
	{
		s.POST("/showtimes", AuthMiddleware(), createShowtime)
		s.GET("/showtimes", listShowtimes)
		s.GET("/showtimes/:id", getShowtime)
		// updates
		s.PUT("/showtimes/:id", AuthMiddleware(), updateShowtime)
		s.DELETE("/showtimes/:id", AuthMiddleware(), deleteShowtime)
	}
}

func createShowtime(c *gin.Context) {
	var body struct{
		MovieID uint `json:"movieId" binding:"required"`
		CinemaID uint `json:"cinemaId" binding:"required"`
		Start string `json:"start" binding:"required"`
		Price float64 `json:"price" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	st, err := time.Parse(time.RFC3339, body.Start)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start time"})
		return
	}
	s := Showtime{MovieID: body.MovieID, CinemaID: body.CinemaID, Start: st, Price: body.Price}
	if err := DB.Create(&s).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

func listShowtimes(c *gin.Context) {
	var ss []Showtime
	DB.Find(&ss)
	c.JSON(http.StatusOK, ss)
}

func getShowtime(c *gin.Context) {
	id := c.Param("id")
	var s Showtime
	if err := DB.First(&s, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func updateShowtime(c *gin.Context) {
	id := c.Param("id")
	var body Showtime
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var s Showtime
	if err := DB.First(&s, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	DB.Model(&s).Updates(body)
	c.JSON(http.StatusOK, s)
}

func deleteShowtime(c *gin.Context) {
	id := c.Param("id")
	DB.Delete(&Showtime{}, id)
	c.Status(http.StatusNoContent)
}
