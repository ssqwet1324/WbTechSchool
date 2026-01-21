package handler

import (
	"net/http"
	"shortener/internal/entity"
	"shortener/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type ShortenerHandler struct {
	uc *usecase.UseCase
}

func New(uc *usecase.UseCase) *ShortenerHandler {
	return &ShortenerHandler{
		uc: uc,
	}
}

func (h *ShortenerHandler) CreateShorten(ctx *ginext.Context) {
	var req entity.ShortenURL

	// парсим json
	if err := ctx.ShouldBind(&req); err != nil {
		zlog.Logger.Error().Err(err).Msg("Error bind json request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// сохраняем и отдаем короткий url
	data, err := h.uc.AddShortUrl(ctx, req)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Error add shorten url")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// RedirectToShorten - перейти по короткой ссылке
func (h *ShortenerHandler) RedirectToShorten(ctx *gin.Context) {
	shortUrl := ctx.Param("short_url")

	// Получаем оригинальный URL
	originalURL, err := h.uc.GetOriginalURL(ctx, shortUrl)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("short_url", shortUrl).Msg("short URL not found")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "short URL not found"})
		return
	}

	// параллельно сохраняем клик в аналитику + проверка на популярный URL
	go func() {
		//TODO сделать чтобы выдавало норм браузер а не длинный
		ua := ctx.Request.UserAgent()

		reqAnalytics := entity.ShortenURLAnalytics{
			ShortURL:    shortUrl,
			TotalClicks: 1,
			ByUserAgent: map[string]int{ua: 1},
			ByDay:       map[string]int{time.Now().Format("2006-01-02"): 1},
			ByMonth:     map[string]int{time.Now().Format("2006-01"): 1},
		}

		// проверяем популярность URL
		if err := h.uc.AddAnalytics(ctx, reqAnalytics); err != nil {
			zlog.Logger.Error().Err(err).Str("short_url", shortUrl).Msg("Failed to add analytics")
		}
	}()

	// редирект на оригинальный URL
	ctx.Redirect(http.StatusTemporaryRedirect, originalURL)

	ctx.JSON(http.StatusTemporaryRedirect, gin.H{"data": originalURL})
}

// GetAnalytics - ручка для аналитики
func (h *ShortenerHandler) GetAnalytics(ctx *gin.Context) {
	shortUrl := ctx.Param("short_url")

	analytics, err := h.uc.GetAnalytics(ctx, shortUrl)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Error get analytics")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"analytics": analytics})
}
