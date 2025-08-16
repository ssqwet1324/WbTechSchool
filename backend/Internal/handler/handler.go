package handler

import (
	"WbDemoProject/Internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	usecase *usecase.Usecase
}

func New(usecase *usecase.Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) GetOrder(ctx *gin.Context) {
	orderUID := ctx.Param("order_uid")

	data, err := h.usecase.GetOrder(ctx, orderUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	ctx.JSON(http.StatusOK, data)
}
