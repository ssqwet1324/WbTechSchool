package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"warehouse_control/internal/entity"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

// Repository - бд
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

// CreateProduct - создание товара
func (repo *Repository) CreateProduct(ctx context.Context, product entity.Product) error {
	query := `INSERT INTO products (product_id, product_name, description, quantity)
    VALUES ($1, $2, $3, $4)`

	_, err := repo.DB.QueryContext(ctx, query, product.ID, product.Name, product.Description, product.Quantity)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: CreateProduct: failed to insert product")
		return err
	}

	return nil
}

// GetProduct - получаем определенный товар по имени
func (repo *Repository) GetProduct(ctx context.Context, productName string) (product entity.Product, err error) {
	query := `SELECT product_id, product_name, description, quantity, updated_at
	FROM products WHERE product_name = $1`

	row := repo.DB.QueryRowContext(ctx, query, productName)
	err = row.Scan(&product.ID, &product.Name, &product.Description, &product.Quantity, &product.UpdatedAt)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: GetProduct: failed to query product")
		return product, err
	}

	return product, nil
}

// GetAllProducts - получаем все товары
func (repo *Repository) GetAllProducts(ctx context.Context) ([]entity.Product, error) {
	query := `SELECT product_id, product_name, description, quantity, updated_at FROM products`

	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: GetAllProducts: failed to query all products")
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("Repository: GetAllProducts: failed to close rows")
		}
	}(rows)

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Quantity, &product.UpdatedAt)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("Repository: GetAllProducts: failed to query all products")
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// UpdateProduct - обновить товар
func (repo *Repository) UpdateProduct(ctx context.Context, product entity.Product) error {
	query := `UPDATE products SET product_name = $1, description = $2, quantity = $3, updated_at = $4
		WHERE product_id = $5`

	result, err := repo.DB.ExecContext(ctx, query, &product.Name, &product.Description, &product.Quantity, &product.UpdatedAt, product.ID)

	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: UpdateProduct: failed to execute update")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: UpdateProduct: failed to get rows affected")
		return err
	}
	if rowsAffected == 0 {
		zlog.Logger.Warn().Msg("Repository: UpdateProduct: no product found with this ID")
	}

	return nil
}

// DeleteProduct - удалить товар
func (repo *Repository) DeleteProduct(ctx context.Context, productName string) error {
	query := `DELETE FROM products WHERE product_name = $1`
	result, err := repo.DB.ExecContext(ctx, query, productName)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: DeleteProduct: failed to execute delete")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: DeleteProduct: failed to get rows affected")
		return err
	}
	if rowsAffected == 0 {
		zlog.Logger.Warn().Msg("Repository: DeleteProduct: no product found with this ID")
	}

	return nil
}

// AddNewUser - добавляем нового пользователя склада
func (repo *Repository) AddNewUser(ctx context.Context, user entity.User) error {
	query := `INSERT INTO users (user_id, username, role) VALUES ($1, $2, $3)`

	_, err := repo.DB.QueryContext(ctx, query, user.ID, user.Username, user.Role)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: AddNewUser: failed to execute insert")
		return err
	}

	zlog.Logger.Info().Msg("Repository: AddNewUser: added user")

	return nil
}

// CheckUser - проверка, существует ли пользователь с таким username
func (repo *Repository) CheckUser(ctx context.Context, username string) (bool, error) {
	query := `SELECT 1 FROM users WHERE username = $1`
	row := repo.DB.QueryRowContext(ctx, query, username)

	var tmp int
	err := row.Scan(&tmp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Пользователя нет
			return false, nil
		}

		zlog.Logger.Error().Err(err).Msg("Repository: CheckUser: failed to query user")
		return false, err
	}

	// Пользователь найден
	return true, nil
}

// GetLogsByProductName - получаем лог для товара
func (repo *Repository) GetLogsByProductName(ctx context.Context, productID string) ([]entity.ProductLogs, error) {
	query := `
		SELECT 
			product_id,
			old_name,
			new_name,
			old_description,
			new_description,
			old_quantity,
			new_quantity,
			changed_at
		FROM product_logs
		WHERE product_id = $1
		ORDER BY changed_at DESC;
	`

	rows, err := repo.DB.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("Repository: GetLogsByProductName: failed to close rows")
		}
	}(rows)

	var logs []entity.ProductLogs
	for rows.Next() {
		var productLogs entity.ProductLogs
		if err := rows.Scan(
			&productLogs.ProductID,
			&productLogs.OldName,
			&productLogs.NewName,
			&productLogs.OldDescription,
			&productLogs.NewDescription,
			&productLogs.OldQuantity,
			&productLogs.NewQuantity,
			&productLogs.ChangedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, productLogs)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

// CheckRole - посмотреть роль
func (repo *Repository) CheckRole(ctx context.Context, username string) (string, error) {
	var role string
	query := `SELECT role FROM users WHERE username = $1`
	row := repo.DB.QueryRowContext(ctx, query, username)

	err := row.Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user %s not found", username)
		}
		zlog.Logger.Error().Err(err).Msg("Repository: CheckRole: failed to query user")
		return "", err
	}

	return role, nil
}
