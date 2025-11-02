package repository

import (
	"context"
	"database/sql"
	"errors"
	"event_booker/internal/entity"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

type Repository struct {
	DB     *dbpg.DB
	master *dbpg.Options
	Client *redis.Client
}

// New - конструктор репозитория
func New(masterDSN string, options *dbpg.Options, client *redis.Client) *Repository {
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
		Client: client,
	}
}

// CreateEvent — создание мероприятия + места
func (repo *Repository) CreateEvent(ctx context.Context, event entity.CreateEvent, totalNumberSeats int, seatNumbers []int) error {
	// проверяем корректное количество мест
	if len(seatNumbers) != totalNumberSeats {
		zlog.Logger.Error().Msg("number of seats does not match total number of seats")
		return fmt.Errorf("seat count mismatch")
	}

	// запускаем транзакцию для безопасного добавления данных в бд
	tx, err := repo.DB.Master.BeginTx(ctx, nil)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to start transaction")
		return err
	}
	defer tx.Rollback()

	// создаем мероприятие
	query := `INSERT INTO events (id, title, date, total_seats, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, query, event.ID, event.Title, event.Date, totalNumberSeats, time.Now())
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("CreateEvent: failed to insert event")
		return err
	}

	// добавляем места
	for _, seatNumber := range seatNumbers {
		seatID := uuid.New().String()
		_, err := tx.ExecContext(ctx,
			`INSERT INTO seats (id, event_id, seat_number, status, user_id) VALUES ($1, $2, $3, $4, $5)`,
			seatID, event.ID, seatNumber, "free", "",
		)
		if err != nil {
			zlog.Logger.Error().Err(err).Msgf("CreateEvent: failed to insert seat %d", seatNumber)
			return err
		}
	}

	return tx.Commit()
}

// CheckFreeSeats — проверка статуса места
func (repo *Repository) CheckFreeSeats(ctx context.Context, eventID string, seatNumber int) (string, error) {
	var status string
	query := `SELECT status FROM seats WHERE event_id=$1 AND seat_number=$2`

	// получаем текущий статус места
	err := repo.DB.QueryRowContext(ctx, query, eventID, seatNumber).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("seat not found")
	}
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("CheckFreeSeats: failed to check seat status")
		return "", err
	}

	return status, nil
}

// MarkSeatAsReserving - помечаем место в бд как зарезервированное
func (repo *Repository) MarkSeatAsReserving(ctx context.Context, eventID string, seatNumber int, userID string) error {
	query := `UPDATE seats SET status='reserving', user_id=$1 WHERE event_id=$2 AND seat_number=$3 
    AND status='free'
    `

	// помечаем место как потенциально зарезервированное
	res, err := repo.DB.ExecContext(ctx, query, userID, eventID, seatNumber)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("MarkSeatAsReserving: failed to update seat")
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("seat %d was already taken", seatNumber)
	}

	return nil
}

// ConfirmSeatBooking — финальное подтверждение после оплаты
func (repo *Repository) ConfirmSeatBooking(ctx context.Context, eventID string, seatNumber int, userID string) error {
	query := `UPDATE seats SET status='booked'
		WHERE event_id=$1 AND seat_number=$2 AND user_id=$3 AND status='reserving'
	`

	// бронируем место после оплаты
	res, err := repo.DB.ExecContext(ctx, query, eventID, seatNumber, userID)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("ConfirmSeatBooking: failed to update seat")
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("reservation not found or already expired")
	}

	return nil
}

// GetEventInfo — получение информации о событии и свободных местах
func (repo *Repository) GetEventInfo(ctx context.Context, eventID string) (*entity.EventInfo, error) {
	// Получаем информацию о событии
	var eventInfo entity.EventInfo
	var date time.Time
	query := `SELECT id, title, date, total_seats FROM events WHERE id=$1`
	err := repo.DB.QueryRowContext(ctx, query, eventID).Scan(&eventInfo.ID, &eventInfo.Title, &date, &eventInfo.TotalSeats)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("event not found")
	}
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("GetEventInfo: failed to get event")
		return nil, err
	}
	eventInfo.Date = date.Format(time.RFC3339)

	// Получаем информацию о местах
	query = `SELECT seat_number, status FROM seats WHERE event_id=$1 ORDER BY seat_number`
	rows, err := repo.DB.QueryContext(ctx, query, eventID)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("GetEventInfo: failed to get seats")
		return nil, err
	}
	defer rows.Close()

	var seats []entity.SeatStatus
	for rows.Next() {
		var seat entity.SeatStatus
		if err := rows.Scan(&seat.SeatNumber, &seat.Status); err != nil {
			zlog.Logger.Error().Err(err).Msg("GetEventInfo: failed to scan seat")
			return nil, err
		}
		seats = append(seats, seat)
	}
	eventInfo.Seats = seats

	return &eventInfo, nil
}

// GetAllEvents — получить все мероприятия из БД
func (repo *Repository) GetAllEvents(ctx context.Context) ([]entity.CreateEvent, error) {
	query := `
        SELECT id, title, date, total_seats
        FROM events
        ORDER BY date ASC;
    `

	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []entity.CreateEvent
	for rows.Next() {
		var e entity.CreateEvent
		if err := rows.Scan(&e.ID, &e.Title, &e.Date, &e.TotalSeats); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return events, nil
}

// CleanupExpiredReservations — очищает просроченные бронирования
//func (repo *Repository) CleanupExpiredReservations(ctx context.Context, redisExist bool) error {
//	// Получаем все места в статусе reserving
//	query := `SELECT id, event_id, seat_number FROM seats WHERE status='reserving'`
//	rows, err := repo.DB.QueryContext(ctx, query)
//	if err != nil {
//		return err
//	}
//	defer rows.Close()
//
//	var toFree []string
//
//	for rows.Next() {
//		var id, eventID string
//		var seatNumber int
//		rows.Scan(&id, &eventID, &seatNumber)
//
//		// если в redis нет ключа
//		if !redisExist {
//			// блокировка отсутствует
//			toFree = append(toFree, id)
//		}
//	}
//
//	// Освобождаем места
//	query = `UPDATE seats SET status='free', user_id='' WHERE id=$1`
//	if len(toFree) > 0 {
//		for _, id := range toFree {
//			_, err := repo.DB.ExecContext(ctx, query, id)
//			if err != nil {
//				zlog.Logger.Error().Err(err).Str("seat_id", id).Msg("failed to free seat")
//			}
//		}
//		zlog.Logger.Info().Msgf("Freed %d expired seats", len(toFree))
//	}
//
//	return nil
//}

func (repo *Repository) CleanupExpiredReservations(ctx context.Context, eventID string, seatNumber int, redisExist bool) error {
	// если ключ есть — ничего не делаем
	if redisExist {
		return nil
	}

	// ключа нет — освобождаем место
	query := `UPDATE seats SET status='free', user_id='' WHERE event_id=$1 AND seat_number=$2 AND status='reserving'`
	_, err := repo.DB.ExecContext(ctx, query, eventID, seatNumber)
	if err != nil {
		zlog.Logger.Error().Err(err).
			Str("event_id", eventID).
			Int("seat_number", seatNumber).
			Msg("failed to free expired seat")
		return err
	}

	zlog.Logger.Info().
		Str("event_id", eventID).
		Int("seat_number", seatNumber).
		Msg("Freed expired seat")

	return nil
}
