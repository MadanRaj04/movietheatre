package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterCinemaRoutes(rg *gin.RouterGroup) {
	s := rg.Group("")
	{
		s.POST("/cinemas", AuthMiddleware(), createCinema)
		s.GET("/cinemas", listCinemas)
		s.GET("/cinemas/:id", getCinema)
		// update/delete with auth
		s.PUT("/cinemas/:id", AuthMiddleware(), updateCinema)
		s.DELETE("/cinemas/:id", AuthMiddleware(), deleteCinema)
	}
}

func createCinema(c *gin.Context) {
	var body Cinema
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := DB.Create(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, body)
}

func listCinemas(c *gin.Context) {
	var cinemas []Cinema
	DB.Find(&cinemas)
	c.JSON(http.StatusOK, cinemas)
}

func getCinema(c *gin.Context) {
	id := c.Param("id")
	var m Cinema
	if err := DB.First(&m, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, m)
}

func updateCinema(c *gin.Context) {
	id := c.Param("id")
	var body Cinema
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var m Cinema
	if err := DB.First(&m, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	DB.Model(&m).Updates(body)
	c.JSON(http.StatusOK, m)
}

func deleteCinema(c *gin.Context) {
	id := c.Param("id")
	DB.Delete(&Cinema{}, id)
	c.Status(http.StatusNoContent)
}
