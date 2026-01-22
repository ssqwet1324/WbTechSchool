package handler

import (
	"fmt"
	"net/http"
	"sales_tracker/internal/entity"
	"sales_tracker/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// Handler - ручки
type Handler struct {
	uc *usecase.UseCase
}

// New - конструктор ручек
func New(uc *usecase.UseCase) *Handler {
	return &Handler{
		uc: uc,
	}
}

// CreateItem - создать запись
func (h *Handler) CreateItem(ctx *ginext.Context) {
	var req entity.NewItem
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Error().Err(err).Msg("handler: CreateItem: Err ShouldBindJSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newItem, err := h.uc.AddItems(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"create_id": newItem})
}

// GetItems - ручка получения записей
func (h *Handler) GetItems(ctx *gin.Context) {
	// Получаем query-параметры
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")

	// Переменные для хранения времени
	fromDate, toDate, err := parseDate(fromStr, toStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Формируем структуру для репозитория
	getItems := entity.GetItems{
		FromDate: fromDate,
		ToDate:   toDate,
	}

	// Вызываем репозиторий
	items, err := h.uc.GetItems(ctx, getItems)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get items"})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

// UpdateItem - ручка обновления записи
func (h *Handler) UpdateItem(ctx *gin.Context) {
	itemID := ctx.Param("id")

	var req entity.NewItem
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Error().Err(err).Msg("handler: UpdateItem: Err ShouldBindJSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.UpdateItems(ctx, req, itemID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"update_id": itemID})
}

// DeleteItem - удалить запись
func (h *Handler) DeleteItem(ctx *gin.Context) {
	itemID := ctx.Param("id")

	err := h.uc.DeleteItems(ctx, itemID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"delete_id": itemID})
}

// GetAnalytics - получить аналитику
func (h *Handler) GetAnalytics(ctx *gin.Context) {
	// Получаем query-параметры
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")
	category := ctx.Query("category")

	var categoryPtr *string
	if category != "" {
		categoryPtr = &category
	}

	// Переменные для хранения времени
	fromDate, toDate, err := parseDate(fromStr, toStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	analytics := entity.GetItemsFromAnalytics{
		FromDate: fromDate,
		ToDate:   toDate,
		Category: categoryPtr,
	}

	data, err := h.uc.GetAnalytics(ctx, analytics)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"analytics": data})
}

// SaveAnalyticsToCSV - сохранить аналитику в csv файл
func (h *Handler) SaveAnalyticsToCSV(ctx *gin.Context) {
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")
	category := ctx.Query("category")
	filename := ctx.Query("filename")

	var categoryPtr *string
	if category != "" {
		categoryPtr = &category
	}

	if filename != "" {
		filename = "analytics.csv"
	}

	// Парсим даты
	fromDate, toDate, err := parseDate(fromStr, toStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	analytics := entity.GetItemsFromAnalytics{
		FromDate: fromDate,
		ToDate:   toDate,
		Category: categoryPtr,
	}

	// Получаем аналитику
	data, err := h.uc.GetAnalytics(ctx, analytics)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Сохраняем CSV
	err = h.uc.SaveAnalyticsToCSV(filename, *data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save CSV: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "CSV saved successfully",
		"filename": filename,
	})
}

// parseDate парсит даты в форматах "YYYY-MM-DD" и RFC3339.
func parseDate(from, to string) (*time.Time, *time.Time, error) {
	var (
		fromDate *time.Time
		toDate   *time.Time
	)

	parse := func(raw string, isFrom bool) (*time.Time, error) {
		if raw == "" {
			return nil, nil
		}

		// Возможные форматы
		layouts := []string{
			time.RFC3339,          // 2025-11-07T14:30:00Z
			"2006-01-02 15:04:05", // 2025-11-07 14:30:00
			"2006-01-02",          // 2025-11-07
		}

		for _, layout := range layouts {
			if t, err := time.Parse(layout, raw); err == nil {
				if layout == "2006-01-02" {
					if isFrom {
						t = t.Truncate(24 * time.Hour) // 00:00:00
					} else {
						t = t.Add(24*time.Hour - time.Second) // 23:59:59
					}
				}
				return &t, nil
			}
		}

		return nil, fmt.Errorf("invalid date format: %s", raw)
	}

	var err error
	if fromDate, err = parse(from, true); err != nil {
		return nil, nil, err
	}
	if toDate, err = parse(to, false); err != nil {
		return nil, nil, err
	}

	return fromDate, toDate, nil
}
