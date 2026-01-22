package handler

import (
	"net/http"
	"time"
	"warehouse_control/internal/entity"
	"warehouse_control/internal/usecase"

	"github.com/gin-gonic/gin"
)

// Handler - ручки
type Handler struct {
	uc *usecase.UseCase
}

// New - конструктор хэндлера
func New(uc *usecase.UseCase) *Handler {
	return &Handler{
		uc: uc,
	}
}

// CreateUser - ручка создания нового пользователя
func (h *Handler) CreateUser(ctx *gin.Context) {
	var user entity.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.AddNewUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"username": user.Username})
}

// Login - войти в аккаунт
func (h *Handler) Login(ctx *gin.Context) {
	var user entity.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.uc.LoginUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"jwt_token": token})
}

// CreateProduct - создать товар
func (h *Handler) CreateProduct(ctx *gin.Context) {
	var product entity.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, err := h.uc.CreateProduct(ctx, product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"product_id": productID})
}

// GetProduct - получить товар
func (h *Handler) GetProduct(ctx *gin.Context) {
	productName := ctx.Param("product_name")

	product, err := h.uc.GetProduct(ctx, productName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"product": product})
}

// GetAllProduct - получение всех товаров
func (h *Handler) GetAllProduct(ctx *gin.Context) {
	products, err := h.uc.GetAllProducts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"products": products})
}

// UpdateProduct - обновить данные о товаре
func (h *Handler) UpdateProduct(ctx *gin.Context) {
	var product entity.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.UpdateProduct(ctx, product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"product_updated": product})
}

// DeleteProduct - удалить товар
func (h *Handler) DeleteProduct(ctx *gin.Context) {
	productName := ctx.Param("product_name")

	err := h.uc.DeleteProduct(ctx, productName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"product_deleted": productName})
}

// GetLogsByProductID - ручка получения истории для товара
func (h *Handler) GetLogsByProductID(ctx *gin.Context) {
	productID := ctx.Param("product_id")
	logs, err := h.uc.GetLogsByProductID(ctx, productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"logs": logs})
}

// SaveProductLogToCSV - ручка сохранения истории изменений в csv
func (h *Handler) SaveProductLogToCSV(ctx *gin.Context) {
	productID := ctx.Param("product_id")
	fileName := productID + "_" + time.Now().Format("2006-01-02_15-04-05") + ".csv"

	// Заголовки для скачивания CSV
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Type", "text/csv")

	// Отправляем CSV напрямую в ответ
	if err := h.uc.SaveHistoryToCSV(ctx, productID, ctx.Writer); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
