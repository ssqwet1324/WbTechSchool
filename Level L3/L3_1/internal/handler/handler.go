package handler

import (
	"L3_1/internal/entity"
	"L3_1/internal/usecase"
	"L3_1/internal/worker"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotifyHandler - конструктор
type NotifyHandler struct {
	useCase *usecase.UseCase
	worker  *worker.Worker
}

// New - конструктор
func New(useCase *usecase.UseCase, worker *worker.Worker) *NotifyHandler {
	return &NotifyHandler{
		useCase: useCase,
		worker:  worker,
	}
}

// CreateNotification - создать уведомление
func (h *NotifyHandler) CreateNotification(ctx *gin.Context) {
	var notify entity.Notify
	if err := ctx.ShouldBindJSON(&notify); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.useCase.CreateNotification(ctx, notify)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.worker.WakeUpWorker()

	ctx.JSON(http.StatusCreated, gin.H{"notification": data})
}

// CheckStatusNotification - посмотреть статус уведомления
func (h *NotifyHandler) CheckStatusNotification(ctx *gin.Context) {
	notifyID := ctx.Param("notifyID")
	status, err := h.useCase.CheckStatusNotification(ctx, notifyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": status})
}

// DeleteNotification - ручка удаления уведомления
func (h *NotifyHandler) DeleteNotification(ctx *gin.Context) {
	notifyID := ctx.Param("notifyID")
	if err := h.useCase.DeleteNotification(ctx, notifyID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.worker.WakeUpWorker()

	ctx.JSON(http.StatusOK, gin.H{"notification deleted successfully:": notifyID})
}

// GetAllNotifications - получить все уведомления
func (h *NotifyHandler) GetAllNotifications(ctx *gin.Context) {
	userID := ctx.Param("userID")
	fmt.Println(userID)
	data, err := h.useCase.GetNotifications(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notifications:": data})
}
