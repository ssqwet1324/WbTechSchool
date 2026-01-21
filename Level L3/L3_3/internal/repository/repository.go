package repository

import (
	"comment_tree/internal/entity"
	"context"
	"database/sql"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

// Repository - структура для работы с бд
type Repository struct {
	DB     *dbpg.DB
	master *dbpg.Options
}

// New - конструктор для репы
func New(masterDSN string, options *dbpg.Options) *Repository {
	masterDB, err := sql.Open("postgres", masterDSN)
	if err != nil {
		zlog.Logger.Fatal().Msgf("failed to open master DB: %v", err)
	}

	masterDB.SetMaxOpenConns(options.MaxOpenConns)
	masterDB.SetMaxIdleConns(options.MaxIdleConns)
	masterDB.SetConnMaxLifetime(options.ConnMaxLifetime)

	db := &dbpg.DB{Master: masterDB}

	return &Repository{DB: db, master: options}
}

// AddComment - добавить комментарий в БД
func (repo *Repository) AddComment(ctx context.Context, newComment entity.NewComment) (*entity.NewComment, error) {
	tx, err := repo.DB.Master.BeginTx(ctx, nil)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: Failed to start transaction")
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// Добавляем сам текст комментария
	flatQuery := `INSERT INTO flat_comments (id, text, created_at) VALUES ($1, $2, now()) RETURNING id, text, created_at`
	var flat entity.NewComment
	err = tx.QueryRowContext(ctx, flatQuery, newComment.CommentID, newComment.CommentText).
		Scan(&flat.CommentID, &flat.CommentText, &flat.CreatedAt)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: Failed to insert into flat_comments")
		return nil, err
	}

	// Добавляем структуру в comments
	commentQuery := `INSERT INTO comments (id, parent_id, comment_ref, created_at) VALUES ($1, $2, $3, now())`
	var parent interface{}
	if newComment.ParentID == nil {
		parent = nil
	} else {
		parent = newComment.ParentID
	}

	_, err = tx.ExecContext(ctx, commentQuery, newComment.CommentID, parent, newComment.CommentID)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: Failed to insert into comments")
		return nil, err
	}

	flat.ParentID = newComment.ParentID

	zlog.Logger.Info().Msgf("Repository: Added new comment %+v", flat.CommentID)

	return &flat, nil
}

// bdQueryGetComments - функция для пробежки получения комментарий
func (repo *Repository) bdQueryGetComments(ctx context.Context, query string, args ...interface{}) (*[]entity.NewComment, error) {
	var err error
	var rows *sql.Rows

	if len(args) == 0 {
		rows, err = repo.DB.QueryContext(ctx, query)
	} else {
		rows, err = repo.DB.QueryContext(ctx, query, args...) // <- распаковываем args
	}

	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: Failed to get comments from DB")
		return nil, err
	}
	defer rows.Close()

	var comments []entity.NewComment
	for rows.Next() {
		var c entity.NewComment
		if err := rows.Scan(&c.CommentID, &c.ParentID, &c.CommentText, &c.CreatedAt); err != nil {
			zlog.Logger.Error().Err(err).Msg("Repository: Failed to scan comments")
			return nil, err
		}
		comments = append(comments, c)
	}

	return &comments, nil
}

// DeleteComment - удалить комментарий (и всех потомков)
func (repo *Repository) DeleteComment(ctx context.Context, commentID string) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := repo.DB.ExecContext(ctx, query, commentID)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: Failed to delete comment")
		return err
	}

	zlog.Logger.Info().Str("commentID", commentID).Msg("Repository: Deleted comment successfully")

	return nil
}

// SearchComment - поиск комментария
func (repo *Repository) SearchComment(ctx context.Context, text string) (*[]entity.NewComment, error) {
	query := `SELECT c.id, c.parent_id, f.text, f.created_at
	FROM comments c JOIN flat_comments f ON c.comment_ref = f.id
	WHERE f.text LIKE '%' || $1 || '%' ORDER BY f.created_at DESC;`

	comments, err := repo.bdQueryGetComments(ctx, query, text)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: Failed to get comments")
		return nil, err
	}

	zlog.Logger.Info().Msg("Repository: SearchComment: comments received successfully")

	return comments, nil
}

// GetParentComments - получить родительские комментарии
func (repo *Repository) GetParentComments(ctx context.Context) (*[]entity.NewComment, error) {
	query := `SELECT c.id, c.parent_id, f.text, c.created_at
    FROM comments c
    JOIN flat_comments f ON c.comment_ref = f.id
    WHERE c.parent_id IS NULL
    ORDER BY c.created_at DESC;`

	comments, err := repo.bdQueryGetComments(ctx, query)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: Failed to get parent comments")
		return nil, err
	}

	return comments, nil
}

// GetChildren - получить детей для одного родителя
func (repo *Repository) GetChildren(ctx context.Context, parentID string, limit, offset int) (*[]entity.NewComment, error) {
	query := `SELECT c.id, c.parent_id, f.text, c.created_at
		FROM comments c
		JOIN flat_comments f ON c.comment_ref = f.id
		WHERE c.parent_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3`

	children, err := repo.bdQueryGetComments(ctx, query, parentID, limit, offset)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: Failed to get children comments")
		return nil, err
	}

	return children, nil
}
