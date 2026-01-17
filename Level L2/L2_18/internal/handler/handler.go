package handler

import (
	"L2_18/internal/entity"
	"L2_18/internal/usecase"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EventHandler - структура для handler
type EventHandler struct {
	Uc  *usecase.UseCase
	Log *zap.Logger
}

// New - конструктор
func New(uc *usecase.UseCase, log *zap.Logger) *EventHandler {
	return &EventHandler{
		Uc:  uc,
		Log: log,
	}
}

// CreateEvent - обрабатываем ручку /create_event
func (eh *EventHandler) CreateEvent(ctx *gin.Context) {
	var event entity.Calendar
	if err := ctx.ShouldBindJSON(&event); err != nil {
		eh.Log.Warn("invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date/time format",
		})
		return
	}

	if err := eh.Uc.SaveEvent(event); err != nil {
		eh.Log.Error("error saving event", zap.Error(err))
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "event created successfully"})
}

// UpdateEvent - обрабатываем ручку /update_event
func (eh *EventHandler) UpdateEvent(ctx *gin.Context) {
	var event entity.Calendar
	if err := ctx.ShouldBindJSON(&event); err != nil {
		eh.Log.Warn("invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date/time format",
		})
		return
	}

	if err := eh.Uc.UpdateEvent(event); err != nil {
		if errors.Is(err, entity.ErrEventNotFound) {
			eh.Log.Warn("event not found", zap.Error(err))
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "event not found"})
			return
		}
		if errors.Is(err, entity.ErrNoEvents) {
			eh.Log.Warn("no events for user", zap.Error(err))
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "no events found"})
			return
		}
		eh.Log.Error("error updating event", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "event updated successfully"})
}

// DeleteEvent - обрабатываем ручку /delete_event
func (eh *EventHandler) DeleteEvent(ctx *gin.Context) {
	var event entity.Calendar

	if err := ctx.ShouldBindJSON(&event); err != nil {
		eh.Log.Warn("invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date/time format",
		})
		return
	}

	if err := eh.Uc.DeleteEvent(event); err != nil {
		if errors.Is(err, entity.ErrEventNotFound) {
			eh.Log.Warn("event not found", zap.Error(err))
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "event not found"})
			return
		}
		if errors.Is(err, entity.ErrNoEvents) {
			eh.Log.Warn("no events for user", zap.Error(err))
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "no events found"})
			return
		}
		eh.Log.Error("error deleting event", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "event deleted successfully"})
}

// EventsForDay - обрабатываем ручку /events_for_day
func (eh *EventHandler) EventsForDay(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	date := ctx.Param("date")

	events, err := eh.Uc.GetEventForDay(userID, date)
	if err != nil {
		if errors.Is(err, entity.ErrParsing) {
			eh.Log.Warn("date parsing error", zap.Error(err), zap.String("date", date))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}
		if errors.Is(err, entity.ErrNoEvents) {
			eh.Log.Info("no events found", zap.String("user_id", userID), zap.String("date", date))
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "no events found for the specified date"})
			return
		}
		eh.Log.Error("error getting events", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": events})
}

// EventsForMonth - обрабатываем ручку /events_for_month
func (eh *EventHandler) EventsForMonth(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	date := ctx.Param("date")

	events, err := eh.Uc.GetEventForMonth(userID, date)
	if err != nil {
		if errors.Is(err, entity.ErrParsing) {
			eh.Log.Warn("date parsing error", zap.Error(err), zap.String("date", date))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}
		if errors.Is(err, entity.ErrNoEvents) {
			eh.Log.Info("no events found", zap.String("user_id", userID), zap.String("date", date))
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "no events found for the specified month"})
			return
		}
		eh.Log.Error("error getting events", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": events})
}

// EventsForWeek - обрабатываем ручку /events_for_week
func (eh *EventHandler) EventsForWeek(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	date := ctx.Param("date")

	events, err := eh.Uc.GetEventsForWeek(userID, date)
	if err != nil {
		if errors.Is(err, entity.ErrParsing) {
			eh.Log.Warn("date parsing error", zap.Error(err), zap.String("date", date))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}
		if errors.Is(err, entity.ErrNoEvents) {
			eh.Log.Info("no events found", zap.String("user_id", userID), zap.String("date", date))
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "no events found for the specified week"})
			return
		}
		eh.Log.Error("error getting events", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": events})
}
