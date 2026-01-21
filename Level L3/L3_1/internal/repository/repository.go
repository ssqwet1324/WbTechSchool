package repository

import (
	"L3_1/internal/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

// Repository - структура для работы с БД
type Repository struct {
	DB     *dbpg.DB
	master *dbpg.Options
}

// New - конструктор для repository
func New(masterDSN string, options *dbpg.Options) *Repository {
	masterDB, err := sql.Open("postgres", masterDSN)
	if err != nil {
		log.Fatalf("failed to open master db: %v", err)
	}

	masterDB.SetMaxOpenConns(options.MaxOpenConns)
	masterDB.SetMaxIdleConns(options.MaxIdleConns)
	masterDB.SetConnMaxLifetime(options.ConnMaxLifetime)

	db := &dbpg.DB{
		Master: masterDB,
		Slaves: nil,
	}

	return &Repository{
		DB:     db,
		master: options,
	}
}

// CreateNotification - добавить уведомление в БД
func (r *Repository) CreateNotification(ctx context.Context, notify entity.Notify) error {
	query := `
		INSERT INTO notifications (user_id, notify_id, title, body, status, sending_date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.DB.ExecContext(ctx, query,
		notify.UserID,
		notify.NotifyID,
		notify.Title,
		notify.Body,
		notify.Status,
		notify.SendingDate,
	)
	if err != nil {
		zlog.Logger.Err(err).Msg("Repository.CreateNotification")
		return fmt.Errorf("error inserting notification: %w", err)
	}

	return nil
}

// CheckStatusNotification - проверка статуса уведомления
func (r *Repository) CheckStatusNotification(ctx context.Context, notifyID uuid.UUID) (string, error) {
	if notifyID.String() == "" {
		return "Notification not found", nil
	}

	query := `SELECT status FROM notifications WHERE notify_id = $1`

	var status bool
	err := r.DB.QueryRowContext(ctx, query, notifyID).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "Notification not found", nil
		}
		zlog.Logger.Err(err).Msg("Repository.CheckStatusNotification")
		return "", fmt.Errorf("error checking status notification: %w", err)
	}

	if status {
		return "Notification sent", nil
	}

	return "Notification not sent", nil
}

// DeleteNotification - удалить уведомление по notifyID
func (r *Repository) DeleteNotification(ctx context.Context, notifyID uuid.UUID) error {
	query := `DELETE FROM notifications WHERE notify_id = $1`
	res, err := r.DB.ExecContext(ctx, query, notifyID)
	if err != nil {
		zlog.Logger.Err(err).Str("notifyID", notifyID.String()).Msg("delete notification query failed")
		return fmt.Errorf("delete notification: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete notification: check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("delete notification: not found (id=%s)", notifyID)
	}

	return nil
}

// GetNotifications - получить все уведомления по userID
func (r *Repository) GetNotifications(ctx context.Context, userID string) ([]entity.Notify, error) {
	query := `SELECT notify_id, title, body, status, sending_date
	FROM notifications WHERE user_id = $1 `
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []entity.Notify
	for rows.Next() {
		var n entity.Notify
		err := rows.Scan(
			&n.NotifyID,
			&n.Title,
			&n.Body,
			&n.Status,
			&n.SendingDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		n.UserID = userID
		notifications = append(notifications, n)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return notifications, nil
}
