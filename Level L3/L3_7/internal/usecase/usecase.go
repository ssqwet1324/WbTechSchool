package usecase

import (
	"context"
	"fmt"
	"io"
	"time"
	"warehouse_control/internal/config"
	"warehouse_control/internal/document"
	"warehouse_control/internal/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/zlog"
)

type RepositoryProvider interface {
	CreateProduct(ctx context.Context, product entity.Product) error
	GetProduct(ctx context.Context, productName string) (product entity.Product, err error)
	GetAllProducts(ctx context.Context) ([]entity.Product, error)
	UpdateProduct(ctx context.Context, product entity.Product) error
	DeleteProduct(ctx context.Context, productName string) error
	AddNewUser(ctx context.Context, user entity.User) error
	CheckUser(ctx context.Context, username string) (bool, error)
	GetLogsByProductName(ctx context.Context, productID string) ([]entity.ProductLogs, error)
	CheckRole(ctx context.Context, username string) (string, error)
}

const (
	accessTokenTTL = time.Hour * 2
)

// UseCase - бизнес логика
type UseCase struct {
	repo RepositoryProvider
	cfg  *config.ServiceConfig
}

// New - конструктор бизнес логики
func New(repo RepositoryProvider, cfg *config.ServiceConfig) *UseCase {
	return &UseCase{
		repo: repo,
		cfg:  cfg,
	}
}

// generateID - сгенерировать новый id
func generateID() string {
	return uuid.New().String()
}

// GenerateJwtToken - создание jwt токена
func (uc *UseCase) GenerateJwtToken(user entity.User, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(accessTokenTTL).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to sign token")
		return "", err
	}

	// Для отладки
	zlog.Logger.Info().Str("token", tokenString).Msg("token generated")

	return tokenString, nil
}

// ExtractUserRoleFromToken - вытаскиваем роль пользователя из jwt токена
func (uc *UseCase) ExtractUserRoleFromToken(tokenStr, secret string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			zlog.Logger.Error().Msgf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		zlog.Logger.Error().Err(err).Msg("ExtractUserRoleFromToken: invalid token")
		return "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		zlog.Logger.Error().Msg("ExtractUserRoleFromToken: cannot parse claims")
		return "", fmt.Errorf("cannot parse claims")
	}

	userRole, ok := claims["role"].(string)
	if !ok || userRole == "" {
		zlog.Logger.Error().Msg("ExtractUserRoleFromToken: userRole not found in token claims")
		return "", fmt.Errorf("userRole not found in token claims")
	}

	return userRole, nil
}

// AddNewUser - создать нового пользователя
func (uc *UseCase) AddNewUser(ctx context.Context, user entity.User) error {
	ok, err := uc.repo.CheckUser(ctx, user.Username)
	if err != nil {
		return err
	}

	if ok {
		zlog.Logger.Warn().Str("username", user.Username).Msg("attempt to add existing user")
		return fmt.Errorf("user %s already exists", user.Username)
	}

	user.ID = generateID()

	return uc.repo.AddNewUser(ctx, user)
}

// LoginUser - войти в аккаунт
func (uc *UseCase) LoginUser(ctx context.Context, user entity.User) (string, error) {
	secretKey := uc.cfg.JWTSecret

	role, err := uc.repo.CheckRole(ctx, user.Username)
	if err != nil {
		return "", err
	}

	if role != user.Role {
		return "", fmt.Errorf("wrong role selected")
	}

	ok, err := uc.repo.CheckUser(ctx, user.Username)
	if err != nil {
		return "", err
	}
	if !ok {
		zlog.Logger.Error().Msg("LoginUser: user not found")
		return "", fmt.Errorf("user not found")
	}

	newUserToken, err := uc.GenerateJwtToken(user, secretKey)
	if err != nil {
		return "", err
	}

	return newUserToken, nil
}

// CreateProduct - создать товар
func (uc *UseCase) CreateProduct(ctx context.Context, product entity.Product) (string, error) {
	if product.ID == "" {
		product.ID = generateID()
	}
	if product.Name == "" {
		zlog.Logger.Error().Msg("CreateProduct: product name is empty")
		return "", fmt.Errorf("product name is empty")
	}

	err := uc.repo.CreateProduct(ctx, product)
	if err != nil {
		return "", err
	}

	return product.ID, nil
}

// GetProduct - получить товар
func (uc *UseCase) GetProduct(ctx context.Context, productName string) (product entity.Product, err error) {
	product, err = uc.repo.GetProduct(ctx, productName)
	if err != nil {
		return product, err
	}

	return product, nil
}

// GetAllProducts - получить все товары на складе
func (uc *UseCase) GetAllProducts(ctx context.Context) ([]entity.Product, error) {
	products, err := uc.repo.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// UpdateProduct - обновить товар
func (uc *UseCase) UpdateProduct(ctx context.Context, product entity.Product) error {
	product.UpdatedAt = time.Now()
	err := uc.repo.UpdateProduct(ctx, product)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProduct - удалить товар со склада
func (uc *UseCase) DeleteProduct(ctx context.Context, productName string) error {
	err := uc.repo.DeleteProduct(ctx, productName)
	if err != nil {
		return err
	}

	return nil
}

// GetLogsByProductID - получить лог изменений в csv
func (uc *UseCase) GetLogsByProductID(ctx context.Context, productID string) ([]entity.ProductLogs, error) {
	if productID == "" {
		zlog.Logger.Error().Msg("GetLogsByProductName: productID is empty")
		return nil, fmt.Errorf("productName is empty")
	}

	logs, err := uc.repo.GetLogsByProductName(ctx, productID)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// SaveHistoryToCSV - сохраняем историю в csv
func (uc *UseCase) SaveHistoryToCSV(ctx context.Context, productID string, w io.Writer) error {
	productHistory, err := uc.GetLogsByProductID(ctx, productID)
	if err != nil {
		return err
	}

	return document.SaveProductHistoryToCSVWriter(w, productHistory)
}
