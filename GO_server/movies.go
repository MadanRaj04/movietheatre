package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterMovieRoutes(rg *gin.RouterGroup) {
	s := rg.Group("")
	{
		s.POST("/movies", AuthMiddleware(), createMovie)
		s.GET("/movies", listMovies)
		s.GET("/movies/:id", getMovie)
		s.PUT("/movies/:id", AuthMiddleware(), updateMovie)
		s.DELETE("/movies/:id", AuthMiddleware(), deleteMovie)
	}
}

func createMovie(c *gin.Context) {
	var body Movie
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

func listMovies(c *gin.Context) {
	var movies []Movie
	DB.Find(&movies)
	c.JSON(http.StatusOK, movies)
}

func getMovie(c *gin.Context) {
	id := c.Param("id")
	var m Movie
	if err := DB.First(&m, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, m)
}

func updateMovie(c *gin.Context) {
	id := c.Param("id")
	var body Movie
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var m Movie
	if err := DB.First(&m, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	DB.Model(&m).Updates(body)
	c.JSON(http.StatusOK, m)
}

func deleteMovie(c *gin.Context) {
	id := c.Param("id")
	DB.Delete(&Movie{}, id)
	c.Status(http.StatusNoContent)
}
