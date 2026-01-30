package handler

import (
	"api_optimization/internal/entity"
	"api_optimization/internal/usecase"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

// Аллоцируем ошибки один раз
var (
	ErrDecode = entity.ErrorResponse{
		Error: entity.ErrorDetail{
			Code: "DecodeError",
		},
	}

	ErrRead = entity.ErrorResponse{
		Error: entity.ErrorDetail{
			Code: "ReadError",
		},
	}
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

// Sum - ручка для обработки суммы
func (h *Handler) Sum(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()
	if err != nil {
		sendError(ctx, ErrRead, err)
		return
	}

	var nums entity.Input

	if err := sonic.Unmarshal(body, &nums); err != nil {
		sendError(ctx, ErrDecode, err)
		return
	}

	sum := h.uc.ReturnSum(nums)

	ctx.JSON(http.StatusOK, entity.Output{Sum: sum.Sum})
}

// sendError - вывод ошибок
func sendError(ctx *gin.Context, base entity.ErrorResponse, err error) {
	resp := base
	resp.Error.Message = err.Error()

	ctx.JSON(http.StatusBadRequest, resp)
}
