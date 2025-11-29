package handler

import (
	"net/http"

	"github.com/bezzang-dev/go-url-shortener/internal/service"
	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(s *service.URLService) *URLHandler {
	return &URLHandler{service: s}
}

type CreateShortURLRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
}

type CreateShortURLResponse struct {
	ShortCode string `json:"short_code"`
}

// CreateShortURL : POST /api/v1/shorten
func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req CreateShortURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	url, err := h.service.GenerateShortURL(req.OriginalURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short URL"})
		return
	}

	c.JSON(http.StatusOK, CreateShortURLResponse{
		ShortCode: url.ShortCode,
	})
}

// RedirectToOriginal : GET /:shortCode
func (h *URLHandler) RedirectToOriginal(c *gin.Context) {
	shortCode := c.Param("shortCode")

	originalURL, err := h.service.GetOriginalURL(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)

}