package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/bezzang-dev/go-url-shortener/internal/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type urlRepository struct {
	db *gorm.DB
	redis *redis.Client
}

func NewURLRepository(db *gorm.DB, rdb *redis.Client) domain.URLRepository {
	return &urlRepository{
		db: db,
		redis: rdb,
	}
}

// GetByShortCode implements domain.URLRepository.
func (u *urlRepository) GetByShortCode(ctx context.Context, code string) (*domain.URL, error) {

	val, err := u.redis.Get(ctx, code).Result()
	if err == redis.Nil {
		// Cache Miss
	} else if err != nil {
		// Redis Error
	} else {
		// Cache Hit
		var url domain.URL
		if err := json.Unmarshal([]byte(val), &url); err == nil {
			return &url, nil
		}
	}

	var url domain.URL
	result := u.db.WithContext(ctx).Where("short_code = ?", code).First(&url)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	jsonBytes, _ := json.Marshal(url)
	u.redis.Set(ctx, code, jsonBytes, 1*time.Hour)
	
	return &url, nil
}

// Save implements domain.URLRepository.
func (u *urlRepository) Save(ctx context.Context, url *domain.URL) error {
	result := u.db.WithContext(ctx).Create(url)
	return result.Error
}