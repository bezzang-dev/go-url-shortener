package repository

import (
	"context"

	"github.com/bezzang-dev/go-url-shortener/internal/domain"
	"gorm.io/gorm"
)

type logRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) domain.AccessLogRepository {
	return &logRepository{
		db: db,
	}
}

func (l *logRepository) Save(ctx context.Context, log *domain.AccessLog) error {
	result := l.db.WithContext(ctx).Create(log)
	return result.Error
}