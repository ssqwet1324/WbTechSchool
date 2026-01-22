package repository

import (
	"context"
	"database/sql"
	"log"
	"sales_tracker/internal/entity"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

// Repository - дб
type Repository struct {
	DB     *dbpg.DB
	master *dbpg.Options
}

// New - конструктор репозитория
func New(masterDSN string, options *dbpg.Options) *Repository {
	masterDB, err := sql.Open("postgres", masterDSN)
	if err != nil {
		log.Fatalf("failed to open master DB: %v", err)
	}

	masterDB.SetMaxOpenConns(options.MaxOpenConns)
	masterDB.SetMaxIdleConns(options.MaxIdleConns)
	masterDB.SetConnMaxLifetime(options.ConnMaxLifetime)
	db := &dbpg.DB{Master: masterDB}

	return &Repository{
		DB:     db,
		master: options,
	}
}

// AddItems - добавить новую запись в бд
func (repo *Repository) AddItems(ctx context.Context, item entity.Item) error {
	query := `INSERT INTO items (id, title, amount, date, category, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := repo.DB.ExecContext(ctx, query, item.ID, item.Title, item.Amount, item.Date, item.Category, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: AddItems: failed to insert item")
		return err
	}

	return nil
}

// GetItems - получить записи
func (repo *Repository) GetItems(ctx context.Context, items entity.GetItems) ([]entity.Item, error) {
	query := `SELECT id, title, amount, date, category, created_at, updated_at FROM items
		WHERE date >= $1 AND date <= $2 ORDER BY date DESC`

	// получаем записи
	rows, err := repo.DB.QueryContext(ctx, query, items.FromDate, items.ToDate)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: GetItems: failed to get items")
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("Repository: GetItems: failed to close rows")
		}
	}(rows)

	// складываем записи в список
	var receivedItems []entity.Item
	for rows.Next() {
		var item entity.Item
		if err := rows.Scan(&item.ID, &item.Title, &item.Amount, &item.Date, &item.Category, &item.CreatedAt,
			&item.UpdatedAt); err != nil {
			zlog.Logger.Error().Err(err).Msg("Repository: GetItems: failed to scan item")
			return nil, err
		}

		receivedItems = append(receivedItems, item)
	}

	return receivedItems, nil
}

// UpdateItems - обновление данных в записи
func (repo *Repository) UpdateItems(ctx context.Context, item entity.Item) error {
	query := `UPDATE items SET title = $1, amount = $2, date = $3, category = $4, updated_at = $5
             WHERE id = $6`
	_, err := repo.DB.ExecContext(ctx, query, item.Title, item.Amount, item.Date, item.Category, item.UpdatedAt, item.ID)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: UpdateItems: failed to update item")
		return err
	}

	return nil
}

// DeleteItems - удаление записи
func (repo *Repository) DeleteItems(ctx context.Context, itemID string) error {
	query := `DELETE FROM items WHERE id = $1`

	_, err := repo.DB.ExecContext(ctx, query, itemID)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: DeleteItems: failed to delete item")
		return err
	}

	return nil
}

// GetAnalytics - получение аналитики
func (repo *Repository) GetAnalytics(ctx context.Context, analytics entity.GetItemsFromAnalytics) (*entity.AnalyticsResult, error) {
	query := `
		WITH filtered AS (
			SELECT amount
			FROM items
			WHERE date >= $1
			  AND date <= $2
			  AND ($3::text IS NULL OR category = $3)
		)
		SELECT
			COUNT(*) AS total_count,
			SUM(amount) AS total_sum,
			AVG(amount) AS avg_amount,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount) AS median_amount,
			PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount) AS p90_amount
		FROM filtered;
		`

	var result entity.AnalyticsResult

	err := repo.DB.QueryRowContext(ctx, query, analytics.FromDate, analytics.ToDate, analytics.Category).
		Scan(&result.TotalCount, &result.TotalSum, &result.AvgAmount, &result.Median, &result.P90)
	if err != nil {
		zlog.Logger.Warn().Err(err).Msg("Repository: GetAnalytics: failed to get items")
		return nil, err
	}

	return &result, nil
}
