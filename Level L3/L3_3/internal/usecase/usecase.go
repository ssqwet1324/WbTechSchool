package usecase

import (
	"comment_tree/internal/entity"
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/zlog"
)

// RepositoryProvider интерфейс репы
type RepositoryProvider interface {
	AddComment(ctx context.Context, newComment entity.NewComment) (*entity.NewComment, error)
	DeleteComment(ctx context.Context, commentID string) error
	SearchComment(ctx context.Context, text string) (*[]entity.NewComment, error)
	GetParentComments(ctx context.Context) (*[]entity.NewComment, error)
	GetChildren(ctx context.Context, parentID string, limit, offset int) (*[]entity.NewComment, error)
}

// UseCase - структура бизнес логики
type UseCase struct {
	repo RepositoryProvider
}

// New конструктор usecase
func New(repo RepositoryProvider) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

// createNewCommentID - создать новый Id
func createNewCommentID() string {
	return uuid.New().String()
}

// AddComment - создать комментарий
func (uc *UseCase) AddComment(ctx context.Context, comment entity.NewComment) (*entity.NewComment, error) {
	newID := createNewCommentID()
	comment.CommentID = newID
	comment.CreatedAt = time.Now()

	data, err := uc.repo.AddComment(ctx, comment)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("UseCase: Failed to add comment")
		return nil, err
	}

	return data, nil
}

// DeleteComment - удалить комментарий
func (uc *UseCase) DeleteComment(ctx context.Context, commentID string) error {
	return uc.repo.DeleteComment(ctx, commentID)
}

// SearchComment - поиск комментария
func (uc *UseCase) SearchComment(ctx context.Context, text string) (*[]entity.NewComment, error) {
	data, err := uc.repo.SearchComment(ctx, text)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("UseCase: Failed to search comment")
		return nil, err
	}

	return data, nil
}

// GetParentComments - получить родительский комент
func (uc *UseCase) GetParentComments(ctx context.Context) (*[]entity.NewComment, error) {
	return uc.repo.GetParentComments(ctx)
}

// GetChildren - получить дочерние комментарии с пагинацией
func (uc *UseCase) GetChildren(ctx context.Context, parentID string, limitStr, offsetStr string) (*[]entity.NewComment, error) {
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("UseCase: Failed to parse limit")
		return nil, err
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("UseCase: Failed to parse offset")
		return nil, err
	}

	return uc.repo.GetChildren(ctx, parentID, limit, offset)
}
