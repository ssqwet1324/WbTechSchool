package entity

import (
	"time"

	"github.com/google/uuid"
)

// Product - товар
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// User - пользователь
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// ProductLogs - логи товаров
type ProductLogs struct {
	ProductID      uuid.UUID `json:"product_id"`
	OldName        string    `json:"old_name"`
	NewName        string    `json:"new_name"`
	OldDescription string    `json:"old_description"`
	NewDescription string    `json:"new_description"`
	OldQuantity    int       `json:"old_quantity"`
	NewQuantity    int       `json:"new_quantity"`
	ChangedAt      time.Time `json:"changed_at"`
}
