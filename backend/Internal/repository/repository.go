package repository

import (
	"WbDemoProject/Internal/config"
	"WbDemoProject/Internal/entity"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB      *pgxpool.Pool
	Config  *config.Config
	Storage map[string]*entity.Order
	Mutex   sync.RWMutex
}

func New(cfg *config.Config) (*Repository, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		strconv.Itoa(cfg.DbPort),
		cfg.DbName,
	)

	//подключаемся к бд
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepository: Error connecrtion from pgxpool: %v", err)
	}

	return &Repository{
		DB:      dbPool,
		Config:  cfg,
		Storage: make(map[string]*entity.Order),
	}, nil
}

// SaveOrderInDB - сохраняем заказ в дб
func (repo *Repository) SaveOrderInDB(ctx context.Context, order *entity.Order) error {
	tx, err := repo.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("PostgresRepository: cannot start transaction: %v", err)
	}

	// при ошибке откатить назад
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			fmt.Printf("PostgresRepository: cannot rollback transaction: %v", err)
		}
	}()

	// orders
	_, err = tx.Exec(ctx, `
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("PostgresRepository: error inserting order: %v", err)
	}

	// delivery
	_, err = tx.Exec(ctx, `INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		return fmt.Errorf("PostgresRepository: error inserting delivery: %v", err)
	}

	// payment
	_, err = tx.Exec(ctx, `
		INSERT INTO payment (order_uid, transaction, request_id, currency,
		                     provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("PostgresRepository: error inserting payment: %v", err)
	}

	// items
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, `
			INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size,
			                   total_price, nm_id, brand, status, order_uid)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
			order.OrderUID,
		)
		if err != nil {
			return fmt.Errorf("PostgresRepository: error inserting item: %v", err)
		}
	}

	//при ошибке транзакции сохраняем данные в кеше чтобы их не потерять
	if err := tx.Commit(ctx); err != nil {
		err := repo.SaveOrderInCache(order)
		if err != nil {
			return fmt.Errorf("repository: cannot commit transaction: %v", err)
		}
		log.Println("заказ не сохранен в бд но сохранен в кэше", err)

		return fmt.Errorf("PostgresRepository: commit error: %v", err)
	}

	log.Printf("заказ сохранен")

	return nil
}

// SaveOrderInCache - сохраняем заказ в кэше
func (repo *Repository) SaveOrderInCache(order *entity.Order) error {
	repo.Mutex.Lock()
	repo.Storage[order.OrderUID] = order
	repo.Mutex.Unlock()

	return nil
}

// GetOrderFromCache - кэшируем и получаем заказ из кэша
func (repo *Repository) GetOrderFromCache(orderUID string) (*entity.Order, bool) {
	repo.Mutex.RLock()
	order, ok := repo.Storage[orderUID]
	repo.Mutex.RUnlock()

	return order, ok
}

// GetOrderFromDB - получаем заказ из бд и кладем в кэш
func (repo *Repository) GetOrderFromDB(ctx context.Context, orderUID string) (*entity.Order, error) {
	var order entity.Order

	err := repo.DB.QueryRow(ctx, `
        SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard 
        FROM orders
        WHERE order_uid = $1`, orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepository: error getting order: %v", err)
	}

	err = repo.DB.QueryRow(ctx, `SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1`, orderUID).Scan(
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepository: error getting delivery: %v", err)
	}

	err = repo.DB.QueryRow(ctx, `SELECT transaction, request_id, currency, provider, amount, payment_dt,
	bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1`, orderUID).Scan(
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDT,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepository: error getting payment: %v", err)
	}

	rows, err := repo.DB.Query(ctx, `SELECT chrt_id, track_number, price, rid, name, sale,
    size, total_price, nm_id, brand, status
	FROM items WHERE order_uid = $1`, orderUID)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepository: error getting items: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Item
		err = rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("PostgresRepository: error getting item: %v", err)
		}

		order.Items = append(order.Items, item)
	}

	return &order, nil
}

// GetAllOrdersFromDB - восстанавливаем кэш
func (repo *Repository) GetAllOrdersFromDB(ctx context.Context) ([]*entity.Order, error) {
	rows, err := repo.DB.Query(ctx, `
       SELECT order_uid, track_number, entry, locale, internal_signature,
              customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
       FROM orders`)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepository: error getting all orders: %v", err)
	}
	defer rows.Close()

	var orders []*entity.Order

	for rows.Next() {
		var order entity.Order
		err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
		)
		if err != nil {
			return nil, fmt.Errorf("PostgresRepository: error scanning order: %v", err)
		}

		// delivery
		err = repo.DB.QueryRow(ctx, `
			SELECT name, phone, zip, city, address, region, email
			FROM delivery WHERE order_uid = $1`, order.OrderUID).Scan(
			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("PostgresRepository: error getting delivery: %v", err)
		}

		// payment
		err = repo.DB.QueryRow(ctx, `
			SELECT transaction, request_id, currency, provider, amount, payment_dt,
				   bank, delivery_cost, goods_total, custom_fee
			FROM payment WHERE order_uid = $1`, order.OrderUID).Scan(
			&order.Payment.Transaction,
			&order.Payment.RequestID,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDT,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
		)
		if err != nil {
			return nil, fmt.Errorf("PostgresRepository: error getting payment: %v", err)
		}

		// items
		itemRows, err := repo.DB.Query(ctx, `
			SELECT chrt_id, track_number, price, rid, name, sale,
			       size, total_price, nm_id, brand, status
			FROM items WHERE order_uid = $1`, order.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("PostgresRepository: error getting items: %v", err)
		}

		for itemRows.Next() {
			var item entity.Item
			err := itemRows.Scan(
				&item.ChrtID,
				&item.TrackNumber,
				&item.Price,
				&item.Rid,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.TotalPrice,
				&item.NmID,
				&item.Brand,
				&item.Status,
			)
			if err != nil {
				itemRows.Close()
				return nil, fmt.Errorf("PostgresRepository: error scanning item: %v", err)
			}
			order.Items = append(order.Items, item)
		}
		itemRows.Close()

		orders = append(orders, &order)
	}

	return orders, nil
}
