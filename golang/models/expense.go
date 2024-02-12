package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseType string

const (
	Food           ExpenseType = "food"
	Utility        ExpenseType = "utility"
	Transportation ExpenseType = "transportation"
	Groceries      ExpenseType = "groceries"
	Subscriptions  ExpenseType = "subscriptions"
	Entertainment  ExpenseType = "entertainment"
	Miscellaneous  ExpenseType = "miscellaneous"
)

type Expense struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Type        ExpenseType `gorm:"type:varchar(20);not null;default:'utility'" json:"type"`
	UserID      uuid.UUID   // Foreign key for User
}

func (base *Expense) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID = uuid.New()
	if base.Type == "" {
		base.Type = Miscellaneous
	}
	return
}
