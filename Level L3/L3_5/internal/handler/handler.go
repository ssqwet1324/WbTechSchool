package handler

import (
	"event_booker/internal/entity"
	"event_booker/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

// EventHandler - ручки
type EventHandler struct {
	uc *usecase.UseCase
}

// New - конструктор
func New(uc *usecase.UseCase) *EventHandler {
	return &EventHandler{
		uc: uc,
	}
}

// CreateEvent - Ручка создания мероприятия
func (h *EventHandler) CreateEvent(ctx *ginext.Context) {
	var newEvent entity.CreateEventRequest

	if err := ctx.ShouldBindJSON(&newEvent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventID, err := h.uc.CreateEvent(ctx, newEvent.Event, newEvent.Layout)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"event_id": eventID})
}

// EventBook бронирование места
func (h *EventHandler) EventBook(ctx *ginext.Context) {
	eventID := ctx.Param("id")

	var req entity.Seat
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.TryReserveSeat(ctx, eventID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "seat reserved successfully", "event_id": eventID, "seat_number": req.SeatNumber})
}

// Confirm - подтверждение оплаты брони
func (h *EventHandler) Confirm(ctx *ginext.Context) {
	eventID := ctx.Param("id")

	var req entity.ConfirmRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.SeatNumber <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "seat_number is required and must be greater than 0"})
		return
	}

	err := h.uc.ConfirmSeatBooking(ctx, eventID, req.SeatNumber, req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "booking confirmed successfully", "event_id": eventID})
}

// GetEvent - получение информации о мероприятии и свободных местах
func (h *EventHandler) GetEvent(ctx *ginext.Context) {
	eventID := ctx.Param("id")

	eventInfo, err := h.uc.GetEvent(ctx, eventID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, eventInfo)
}

// GetAllEvents - получение списка всех мероприятий
func (h *EventHandler) GetAllEvents(ctx *ginext.Context) {
	events, err := h.uc.GetAllEvents(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, events)
}
