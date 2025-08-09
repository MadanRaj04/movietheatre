package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	godotenv.Load()

}

var DB *gorm.DB

func main() {
	dsn := os.Getenv("MYSQL_DSN") // e.g. user:pass@tcp(127.0.0.1:3306)/dbname?parseTime=true
	if dsn == "" {
		dsn = "root:honda948@tcp(127.0.0.1:3306)/cinema?parseTime=true"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	DB = db

	// Auto-migrate models
	if err := DB.AutoMigrate(&User{}, &Movie{}, &Cinema{}, &Showtime{}, &Reservation{}); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	router := gin.Default()
	router.Static("/uploads", "./uploads")
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	api := router.Group("")
	{
		RegisterUserRoutes(api)
		RegisterMovieRoutes(api)
		RegisterCinemaRoutes(api)
		RegisterShowtimeRoutes(api)
		RegisterReservationRoutes(api)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	s := &httpServer{engine: router, addr: ":" + port}
	s.run()
}

type httpServer struct {
	engine *gin.Engine
	addr   string
}

func (s *httpServer) run() {
	log.Printf("listening %s", s.addr)
	if err := s.engine.Run(s.addr); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
