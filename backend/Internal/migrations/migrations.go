package migrations

import (
	"WbDemoProject/Internal/repository"
	"context"
	"fmt"
	"time"
)

type Migrations struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Migrations {
	return &Migrations{repo: repo}
}

func (m *Migrations) InitTables(ctx context.Context) error {
	query := `
CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(50) PRIMARY KEY,
    track_number VARCHAR(50),
    entry VARCHAR(20),
    locale VARCHAR(10),
    internal_signature VARCHAR(100),
    customer_id VARCHAR(50),
    delivery_service VARCHAR(50),
    shardkey VARCHAR(10),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(10)
);

CREATE TABLE IF NOT EXISTS delivery (
    order_uid VARCHAR(50) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR(100),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(50),
    address VARCHAR(100),
    region VARCHAR(50),
    email VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS payment (
    order_uid VARCHAR(50) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction VARCHAR(50),
    request_id VARCHAR(50),
    currency VARCHAR(10),
    provider VARCHAR(50),
    amount INT,
    payment_dt BIGINT,
    bank VARCHAR(50),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id BIGINT,
    track_number VARCHAR(50),
    price INT,
    rid VARCHAR(50),
    name VARCHAR(100),
    sale INT,
    size VARCHAR(10),
    total_price INT,
    nm_id BIGINT,
    brand VARCHAR(50),
    status INT
);
`
	maxRetries := 5
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		_, err := m.repo.DB.Exec(ctx, query)
		if err == nil {
			// Всё ок, таблицы созданы
			return nil
		}

		fmt.Printf("Не удалось создать таблицы (попытка %d/%d): %v\n", i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("не удалось создать таблицы после %d попыток", maxRetries)
}
