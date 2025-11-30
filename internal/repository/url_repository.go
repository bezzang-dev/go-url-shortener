package repository

import (
	"errors"

	"github.com/bezzang-dev/go-url-shortener/internal/domain"
	"gorm.io/gorm"
)

type urlRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) domain.URLRepository {
	return &urlRepository{
		db: db,
	}
}

// GetByShortCode implements domain.URLRepository.
func (u *urlRepository) GetByShortCode(code string) (*domain.URL, error) {
	var url domain.URL

	result := u.db.Where("short_code = ?", code).First(&url)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &url, nil
}

// Save implements domain.URLRepository.
func (u *urlRepository) Save(url *domain.URL) error {
	result := u.db.Create(&url)

	if result.Error != nil {
		return result.Error
	}
	return nil
}