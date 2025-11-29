package service

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/bezzang-dev/go-url-shortener/internal/domain"
)

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type URLService struct {
	repo domain.URLRepository
}

func NewURLService(repo domain.URLRepository) *URLService {
	return &URLService{
		repo: repo,
	}
}

func (s *URLService) GenerateShortURL(originalURL string) (*domain.URL, error) {

	shortCode, err := generateRandomString(6)
	if err != nil {
		return nil, err
	}

	url := &domain.URL{
		Original: originalURL,
		ShortCode: shortCode,
	}

	if err := s.repo.Save(url); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *URLService) GetOriginalURL(shortCode string) (string, error) {
	url, err := s.repo.GetByShortCode(shortCode)
	if err != nil {
		return "", err
	}
	return url.Original, nil
}

func generateRandomString(length int) (string, error) {
	var sb strings.Builder
	for range length {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(base62Chars))))
		if err != nil {
			return "", err
		}
		sb.WriteByte(base62Chars[num.Int64()])
	}
	return sb.String(), nil
}