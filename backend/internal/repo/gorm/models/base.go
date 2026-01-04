package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"primarykey;size:36"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
