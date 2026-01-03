package domain

import (
	"context"
	"time"
)

type AccessLog struct {
	ID uint `gorm:"primaryKey" json:"-"`
	ShortCode string `gorm:"size:10;index:idx_short_code;not null" json:"short_code"`
	ClientIP string `gorm:"size:45;not null" json:"client_ip"`
	UserAgent string `gorm:"size:255;not null" json:"user_agent"`
	AccessedAt time.Time `gorm:"not null" json:"accessed_at"`
}

type AccessLogRepository interface {
	Save(ctx context.Context, accessLog *AccessLog) error
}