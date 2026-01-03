package main

import (
	"log"

	"github.com/bezzang-dev/go-url-shortener/internal/analytics"
	"github.com/bezzang-dev/go-url-shortener/internal/domain"
	"github.com/bezzang-dev/go-url-shortener/internal/handler"
	"github.com/bezzang-dev/go-url-shortener/internal/repository"
	"github.com/bezzang-dev/go-url-shortener/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	dsn := "host=localhost user=postgres password=postgres dbname=shortener port=5432 sslmode=disable TimeZone=Asia/Seoul"

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&domain.URL{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	analyticsClient, err := analytics.NewClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect to analytics service: %v", err)
	}
	defer analyticsClient.Close()

	repo := repository.NewURLRepository(db, rdb)
	svc := service.NewURLService(repo)
	h := handler.NewURLHandler(svc, analyticsClient)

	r := gin.Default()

	r.SetTrustedProxies(nil)

	api := r.Group("/api/v1") 
	{
		api.POST("/shorten", h.CreateShortURL)
	}

	r.GET("/:shortCode", h.RedirectToOriginal)

	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}

}