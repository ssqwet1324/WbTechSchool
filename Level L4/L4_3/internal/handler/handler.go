package handler

import (
	"L4_3/internal/entity"
	"L4_3/internal/log"
	"L4_3/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EventHandler - структура для handler
type EventHandler struct {
	uc  *usecase.UseCase
	log *log.Log
}

// New - конструктор
func New(uc *usecase.UseCase, log *log.Log) *EventHandler {
	return &EventHandler{
		uc:  uc,
		log: log,
	}
}

// CreateEvent - обрабатываем ручку /create_event
func (eh *EventHandler) CreateEvent(ctx *gin.Context) {
	var event entity.Calendar
	if err := ctx.ShouldBindJSON(&event); err != nil {
		eh.log.AsyncError("error parsing json", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := eh.uc.SaveEvent(event); err != nil {
		eh.log.AsyncError("error saving event", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Created Event success": event})
}

// UpdateEvent - обрабатываем ручку /update_event
func (eh *EventHandler) UpdateEvent(ctx *gin.Context) {
	var event entity.Calendar
	if err := ctx.ShouldBindJSON(&event); err != nil {
		eh.log.AsyncError("error parsing json", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := eh.uc.UpdateEvent(event); err != nil {
		eh.log.AsyncError("error updating event", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Updated Event success": event})
}

// DeleteEvent - обрабатываем ручку /delete_event
func (eh *EventHandler) DeleteEvent(ctx *gin.Context) {
	var event entity.Calendar

	if err := ctx.ShouldBindJSON(&event); err != nil {
		eh.log.AsyncError("error parsing json", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := eh.uc.DeleteEvent(event); err != nil {
		eh.log.AsyncError("error deleting event", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Deleted Event success": event.DataEvent})
}

// EventsForDay - обрабатываем ручку /events_for_day
func (eh *EventHandler) EventsForDay(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	date := ctx.Param("date")

	event, err := eh.uc.GetEventForDay(userID, date)
	if err != nil {
		eh.log.AsyncError("error getting event", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"event": event})
}

// EventsForMonth - обрабатываем ручку /events_for_month
func (eh *EventHandler) EventsForMonth(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	date := ctx.Param("date")
	event, err := eh.uc.GetEventForMonth(userID, date)
	if err != nil {
		eh.log.AsyncError("error getting event", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"event": event})
}

// EventsForWeek - обрабатываем ручку /events_for_week
func (eh *EventHandler) EventsForWeek(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	date := ctx.Param("date")
	event, err := eh.uc.GetEventsForWeek(userID, date)
	if err != nil {
		eh.log.AsyncError("error getting events", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"events": event})
}
