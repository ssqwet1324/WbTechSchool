package handler

import (
	"encoding/json"
	"image_processor/internal/entity"
	"image_processor/internal/kafka"
	"image_processor/internal/usecase"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

// ImgHandler - структура handler
type ImgHandler struct {
	uc *usecase.UseCase
	pr *kafka.Queue
}

// New - конструктор для ImgHandler
func New(uc *usecase.UseCase, pr *kafka.Queue) *ImgHandler {
	return &ImgHandler{
		uc: uc,
		pr: pr,
	}
}

// UploadImage - ручка загрузки фото
func (h *ImgHandler) UploadImage(ctx *ginext.Context) {
	// получаем изображение
	img, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
		return
	}

	defer func(img multipart.File) {
		err := img.Close()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close file"})
		}
	}(img)

	imgName := header.Filename
	imgSize := header.Size

	loadPhoto := entity.LoadPhoto{
		Name:      imgName,
		Size:      imgSize,
		Reader:    img,
		Status:    "loading",
		CreatedAt: time.Now(),
	}

	// добавляем фото
	photoID, err := h.uc.AddPhoto(ctx, loadPhoto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add photo"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"photo_id": photoID})
}

// PhotoProcessing - обработка фотографии
func (h *ImgHandler) PhotoProcessing(ctx *gin.Context) {
	var req entity.PhotoInfo

	// Получаем параметры размера из query параметров для resize и miniature
	if req.Version == "resize" || req.Version == "miniature" {
		widthPhoto := ctx.Query("widthPhoto")
		heightPhoto := ctx.Query("heightPhoto")

		req.Width = widthPhoto
		req.Height = heightPhoto
	}

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем сообщение для кафки
	msg, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal photo info"})
		return
	}

	// отправляем в Kafka
	if err := h.pr.SendMessage(ctx.Request.Context(), []byte(req.PhotoID), msg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue photo for processing"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"start_processing": req.PhotoID})
}

// GetProcessedImg - получить обработанное изображение
func (h *ImgHandler) GetProcessedImg(ctx *gin.Context) {
	var req entity.PhotoInfo
	photoVersion := ctx.Param("photo_version")
	photoID := ctx.Param("id")

	req.PhotoID = photoID
	req.Version = photoVersion

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imgUrl, err := h.uc.GetProcessedImg(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"img_url": imgUrl})
}

// DeletePhoto - ручка удаления фотографии
func (h *ImgHandler) DeletePhoto(ctx *gin.Context) {
	var req entity.PhotoInfo

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.DeletePhoto(ctx, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"deleted_photo": req.PhotoID})
}
