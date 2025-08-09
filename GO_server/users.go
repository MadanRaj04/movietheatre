package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	s := rg.Group("")
	{
		s.POST("/users", registerHandler)
		s.POST("/users/login", loginHandler)
		// protected
		p := s.Group("")
		p.Use(AuthMiddleware())
		p.GET("/me", meHandler)
	}
}

func registerHandler(c *gin.Context) {
	var body struct {
		Name     string `json:"Name" binding:"required"`
		Username string `json:"Username" binding:"required"`
		Email    string `json:"Email" binding:"required,email"`
		Password string `json:"Password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// hash
	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	user := User{Name: body.Name, Email: body.Email, Password: string(hash), Role: "user", Username: body.Username}
	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func loginHandler(c *gin.Context) {
	var body struct {
		Email    string `json:"Email" binding:"required,email"`
		Password string `json:"Password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user User
	if err := DB.Where("Email = ?", body.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	// token
	token, err := GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func meHandler(c *gin.Context) {
	uid := c.GetUint("userId")
	var user User
	if err := DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"user": user})
}
