package entities

import "github.com/google/uuid"

type Crypto struct {
	ID    uuid.UUID `gorm:"primaryKey" json:"id"`
	Asset string    `gorm:"not null" json:"asset"`
	Value int       `json:"value"`
	User  *User     `json:"user,omitempty"`
}
