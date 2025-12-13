package service

import (
	"context"
	"testing"

	"github.com/bezzang-dev/go-url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) Save(ctx context.Context, url *domain.URL) error {
	args := m.Called(url)
	return args.Error(0)
}

func (m *MockURLRepository) GetByShortCode(ctx context.Context, code string) (*domain.URL, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.URL), args.Error(1)
}

func TestGenerateShortUrl(t *testing.T) {
	mockRepo := new(MockURLRepository)
	svc := NewURLService(mockRepo)
	originalURL := "https://www.naver.com"

	ctx := context.Background()

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	// when
	result, err := svc.GenerateShortURL(ctx, originalURL)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, originalURL, result.Original)
	assert.Len(t, result.ShortCode, 6)

	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	svc := NewURLService(mockRepo)

	shortCode := "AbC12z"
	expectedURL := "https://www.google.com"

	ctx := context.Background()

	mockURL := &domain.URL{
		Original: expectedURL,
		ShortCode: shortCode,
	}

	mockRepo.On("GetByShortCode", shortCode).Return(mockURL, nil)

	resultStr, err := svc.GetOriginalURL(ctx, shortCode)

	assert.NoError(t, err)
	assert.Equal(t, expectedURL, resultStr)

	mockRepo.AssertExpectations(t)
}