package handler

import (
	"comment_tree/internal/entity"
	"comment_tree/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// CommentHandler хендлер
type CommentHandler struct {
	uc *usecase.UseCase
}

// New конструктор хендлера
func New(uc *usecase.UseCase) *CommentHandler {
	return &CommentHandler{
		uc: uc,
	}
}

// AddComment - ручка создания комментов
func (h *CommentHandler) AddComment(ctx *ginext.Context) {
	var comment entity.NewComment

	if err := ctx.ShouldBindJSON(&comment); err != nil {
		zlog.Logger.Error().Err(err).Msg("Error binding JSON request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.uc.AddComment(ctx, comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"new_comment": data})
}

// DeleteComment - удалить комменты
func (h *CommentHandler) DeleteComment(ctx *ginext.Context) {
	commentID := ctx.Param("comment_id")

	err := h.uc.DeleteComment(ctx, commentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"deleted": commentID})
}

// SearchComment - ручка поиска комментариев
func (h *CommentHandler) SearchComment(ctx *ginext.Context) {
	var searchComment entity.SearchComment
	if err := ctx.ShouldBindJSON(&searchComment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.uc.SearchComment(ctx, searchComment.Text)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"searched_comments": data})
}

// GetParentComments - получить родительский коммент
func (h *CommentHandler) GetParentComments(ctx *ginext.Context) {
	data, err := h.uc.GetParentComments(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"parent_comments": data})
}

// GetComments - получить дочерние комментарии для одного родителя
func (h *CommentHandler) GetComments(ctx *ginext.Context) {
	parentID := ctx.Query("parent_id")
	if parentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parent_id is required"})
		return
	}

	limitStr := ctx.DefaultQuery("limit", "20")
	offsetStr := ctx.DefaultQuery("offset", "0")

	data, err := h.uc.GetChildren(ctx, parentID, limitStr, offsetStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"children": data})
}
