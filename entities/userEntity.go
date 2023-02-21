package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;type:varchar(255)" json:"email" `
	Password  string    `gorm:"->;<-;not null" json:"-" validate:"required, min=6"`
	Token     string    `gorm:"-" json:"token,omitempty"`
	Asset     string    `json:"asset"`
	Value     int       `json:"value"`
	CryptoID  string    `gorm:"default:null" json:"crypto_id"`
	Crypto    Crypto    `gorm:"foreignkey:CryptoID;constraint:onUpdate:CASCADE,onDelete:SET NULL" json:"crypto"`
	PostID    string    `gorm:"default:null" json:"post_id"`
	Posts     []*Post   `json:"posts,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
