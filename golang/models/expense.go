package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Expense struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Description string
	Amount      float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      uuid.UUID // Foreign key for User
}

func (base *Expense) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID = uuid.New()
	return
}
