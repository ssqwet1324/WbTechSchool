package entity

import "time"

// CreateEvent - создание мероприятия
type CreateEvent struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Date       time.Time `json:"date"`
	TotalSeats int       `json:"total_seats"`
}

// TotalSeats - количество сидений
type TotalSeats struct {
	Rows        int `json:"rows"`
	SeatsPerRow int `json:"seats_per_row"`
	StartNumber int `json:"start_number"`
}

// Seat - номер сидения
type Seat struct {
	SeatNumber int    `json:"seat_number"`
	SeatID     string `json:"seat_id"`
	UserID     string `json:"user_id"`
}

// CreateEventRequest - объедененный json
type CreateEventRequest struct {
	Event  CreateEvent `json:"event"`
	Layout TotalSeats  `json:"layout"`
}

// EventInfo - информация о событии
type EventInfo struct {
	ID         string       `json:"id"`
	Title      string       `json:"title"`
	Date       string       `json:"date"`
	TotalSeats int          `json:"total_seats"`
	Seats      []SeatStatus `json:"seats"`
}

// SeatStatus - статус места
type SeatStatus struct {
	SeatNumber int    `json:"seat_number"`
	Status     string `json:"status"`
}

// ConfirmRequest - запрос на подтверждение оплаты
type ConfirmRequest struct {
	UserID     string `json:"user_id"`
	SeatNumber int    `json:"seat_number"`
}
