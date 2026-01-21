package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// UseCase Interface
type UseCaseInterface interface {
	AddShortURL(ctx *gin.Context, req AddShortURLRequest) (string, error)
	GetOriginalURL(ctx *gin.Context, shortURL string) (string, error)
	AddAnalytics(ctx *gin.Context, analytics ShortURLAnalytics) error
	GetAnalytics(ctx *gin.Context, shortURL string) (ShortURLAnalytics, error)
}

// Mock UseCase
type mockUseCase struct{}

func (m *mockUseCase) AddShortURL(ctx *gin.Context, req AddShortURLRequest) (string, error) {
	return "test123", nil
}

func (m *mockUseCase) GetOriginalURL(ctx *gin.Context, shortURL string) (string, error) {
	if shortURL == "notfound" {
		return "", ErrNotFound
	}
	return "https://example.com", nil
}

func (m *mockUseCase) AddAnalytics(ctx *gin.Context, analytics ShortURLAnalytics) error {
	return nil
}

func (m *mockUseCase) GetAnalytics(ctx *gin.Context, shortURL string) (ShortURLAnalytics, error) {
	if shortURL == "notfound" {
		return ShortURLAnalytics{}, ErrNotFound
	}
	return ShortURLAnalytics{
		ShortURL:    shortURL,
		TotalClicks: 10,
		ByDay:       map[string]int{"2025-10-07": 5},
		ByMonth:     map[string]int{"2025-10": 10},
		ByUserAgent: map[string]int{"Chrome": 7},
	}, nil
}

// Entities
type AddShortURLRequest struct {
	OriginalURL string `json:"original_url"`
}

type ShortURLAnalytics struct {
	ShortURL    string         `json:"short_url"`
	TotalClicks int            `json:"total_clicks"`
	ByDay       map[string]int `json:"by_day"`
	ByMonth     map[string]int `json:"by_month"`
	ByUserAgent map[string]int `json:"by_user_agent"`
}

var ErrNotFound = http.ErrMissingFile

// Handler
type TestShortenerHandler struct {
	uc UseCaseInterface
}

func NewHandler(uc UseCaseInterface) *TestShortenerHandler {
	return &TestShortenerHandler{uc: uc}
}

func (h *TestShortenerHandler) CreateShorten(ctx *gin.Context) {
	var req AddShortURLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortURL, err := h.uc.AddShortURL(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": shortURL})
}

func (h *TestShortenerHandler) RedirectToShorten(ctx *gin.Context) {
	shortURL := ctx.Param("short_url")

	originalURL, err := h.uc.GetOriginalURL(ctx, shortURL)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "short URL not found"})
		return
	}

	go func() {
		ua := ctx.Request.UserAgent()
		analytics := ShortURLAnalytics{
			ShortURL:    shortURL,
			TotalClicks: 1,
			ByDay:       map[string]int{time.Now().Format("2006-01-02"): 1},
			ByMonth:     map[string]int{time.Now().Format("2006-01"): 1},
			ByUserAgent: map[string]int{ua: 1},
		}
		_ = h.uc.AddAnalytics(ctx, analytics)
	}()

	ctx.Redirect(http.StatusTemporaryRedirect, originalURL)
}

func (h *TestShortenerHandler) GetAnalytics(ctx *gin.Context) {
	shortURL := ctx.Param("short_url")

	analytics, err := h.uc.GetAnalytics(ctx, shortURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"analytics": analytics})
}

// Setup Router
func setupTestRouter(handler *TestShortenerHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/shorten", handler.CreateShorten)
	router.GET("/s/:short_url", handler.RedirectToShorten)
	router.GET("/analytics/:short_url", handler.GetAnalytics)

	return router
}

// Tests
func TestHandler_CreateShorten(t *testing.T) {
	mockUC := &mockUseCase{}
	handler := NewHandler(mockUC)
	router := setupTestRouter(handler)

	req := httptest.NewRequest("POST", "/shorten", strings.NewReader(`{"original_url":"https://example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["data"] != "test123" {
		t.Errorf("expected data=test123, got %v", resp["data"])
	}
}

func TestHandler_RedirectToShorten(t *testing.T) {
	mockUC := &mockUseCase{}
	handler := NewHandler(mockUC)
	router := setupTestRouter(handler)

	req := httptest.NewRequest("GET", "/s/test123", nil)
	req.Header.Set("User-Agent", "TestAgent")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected status 307, got %d", w.Code)
	}

	location := w.Header().Get("Location")
	if location != "https://example.com" {
		t.Errorf("expected redirect to https://example.com, got %s", location)
	}
}

func TestHandler_GetAnalytics(t *testing.T) {
	mockUC := &mockUseCase{}
	handler := NewHandler(mockUC)
	router := setupTestRouter(handler)

	req := httptest.NewRequest("GET", "/analytics/test123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["analytics"] == nil {
		t.Error("expected analytics field in response")
	}
}
