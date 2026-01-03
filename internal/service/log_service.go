package service

import (
	"context"
	"time"

	"github.com/bezzang-dev/go-url-shortener/internal/domain"
)
type LogService struct {
	repo domain.AccessLogRepository
}

func NewLogService(repo domain.AccessLogRepository) *LogService {
	return &LogService{
		repo: repo,
	}
}

func (s *LogService) RecordAccess(context context.Context, shortCode, ip, userAgent string, accessedAt int64) error {
	log := &domain.AccessLog{
		ShortCode:  shortCode,
		ClientIP:   ip,
		UserAgent:  userAgent,
		AccessedAt: time.Unix(accessedAt, 0),
	}
	return s.repo.Save(context, log)
}