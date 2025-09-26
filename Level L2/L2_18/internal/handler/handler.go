package handler

import (
	"L2_18/internal/entity"
	"L2_18/internal/usecase"
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
		eh.Log.Error("error parsing json", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := eh.Uc.SaveEvent(event); err != nil {
		eh.Log.Error("error saving event", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Created Event success": event})
}

// UpdateEvent - обрабатываем ручку /update_event
func (eh *EventHandler) UpdateEvent(ctx *gin.Context) {
	var event entity.Calendar
	if err := ctx.ShouldBindJSON(&event); err != nil {
		eh.Log.Error("error parsing json", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := eh.Uc.UpdateEvent(event); err != nil {
		eh.Log.Error("error updating event", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Updated Event success": event})
}

// DeleteEvent - обрабатываем ручку /delete_event
func (eh *EventHandler) DeleteEvent(ctx *gin.Context) {
	var event entity.Calendar

	if err := ctx.ShouldBindJSON(&event); err != nil {
		eh.Log.Error("error parsing json", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := eh.Uc.DeleteEvent(event); err != nil {
		eh.Log.Error("error deleting event", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Deleted Event success": event.DataEvent})
}

// EventsForDay - обрабатываем ручку /events_for_day
func (eh *EventHandler) EventsForDay(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	date := ctx.Param("date")

	event, err := eh.Uc.GetEventForDay(userID, date)
	if err != nil {
		eh.Log.Error("error getting event", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"event": event})
}

// EventsForMonth - обрабатываем ручку /events_for_month
func (eh *EventHandler) EventsForMonth(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	date := ctx.Param("date")
	event, err := eh.Uc.GetEventForMonth(userID, date)
	if err != nil {
		eh.Log.Error("error getting event", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"event": event})
}

// EventsForWeek - обрабатываем ручку /events_for_week
func (eh *EventHandler) EventsForWeek(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	date := ctx.Param("date")
	event, err := eh.Uc.GetEventsForWeek(userID, date)
	if err != nil {
		eh.Log.Error("error getting events", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"events": event})
}
