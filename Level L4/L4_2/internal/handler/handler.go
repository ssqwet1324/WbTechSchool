package handler

import (
	"net/http"

	"L4.2/internal/entity"
	"L4.2/internal/search"
	"github.com/gin-gonic/gin"
)

type WorkerHandler struct{}

func NewWorkerHandler() *WorkerHandler {
	return &WorkerHandler{}
}

func (h *WorkerHandler) HandleSearch(ctx *gin.Context) {
	var req entity.SearchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := search.SearchLines(req.Lines, req.Options, req.Pattern, req.Offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entity.SearchResponse{
		Lines: result.Lines,
		Count: result.Count,
	})
}
