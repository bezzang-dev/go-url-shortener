package domain

import (
	"context"
	"time"
)

type URL struct {
	ID uint `gorm:"primaryKey" json:"-"`
	Original string `gorm:"type:text;not null" json:"original_url"`
	ShortCode string `gorm:"size:10;uniqueIndex;not null" json:"short_code"`
	CreatedAt time.Time `json:"created_at"`
}

type URLRepository interface {
	Save(ctx context.Context, url *URL) error
	GetByShortCode(ctx context.Context, code string) (*URL, error)
}