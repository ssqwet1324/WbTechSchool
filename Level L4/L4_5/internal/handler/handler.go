package handler

import (
	"api_optimization/internal/entity"
	"api_optimization/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler - ручки
type Handler struct {
	uc *usecase.UseCase
}

// New - конструктор
func New(uc *usecase.UseCase) *Handler {
	return &Handler{
		uc: uc,
	}
}

// Sum - ручка обработки суммы
func (h *Handler) Sum(ctx *gin.Context) {
	var nums entity.Input
	if err := ctx.ShouldBindJSON(&nums); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sum := h.uc.ReturnSum(nums)

	ctx.JSON(http.StatusOK, entity.Output{Sum: sum.Sum})
}
