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

	// Сначала ищем в кэше
	data, exist := h.usecase.GetOrderFromCache(orderUID)
	if !exist {
		var err error
		data, err = h.usecase.GetOrderFromDB(ctx, orderUID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})

			return
		}
		_ = h.usecase.SaveOrderInDB(ctx, data)
	}

	ctx.JSON(http.StatusOK, data)
}
