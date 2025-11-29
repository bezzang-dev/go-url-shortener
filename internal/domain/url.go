package domain

import "time"

type URL struct {
	ID uint `gorm:"primaryKey" json:"-"`
	Original string `gorm:"type:text;not null" json:"original_url"`
	ShortCode string `gorm:"size:10;uniqueIndex;not null" json:"short_code"`
	CreatedAt time.Time `json:"created_at"`
}

type URLRepository interface {
	Save(url *URL) error
	GetByShortCode(code string) (*URL, error)
}