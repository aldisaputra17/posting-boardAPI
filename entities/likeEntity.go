package entities

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	Likes     int       `gorm:"default:0" json:"likes"`
	PostID    string    `gorm:"not null" json:"post_id"`
	Post      *Post     `gorm:"foreignkey:PostID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"-"`
	UserID    string    `gorm:"not null" json:"user_id"`
	User      *User     `gorm:"foreignkey:UserID;constraint:onUpdate:CASCADE,onDelete:SET NULL" json:"user"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
