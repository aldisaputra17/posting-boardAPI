package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID  `gorm:"primaryKey" json:"id"`
	Asset       string     `gorm:"not null" json:"asset"`
	Fiat        string     `gorm:"not null" json:"fiat"`
	Description string     `gorm:"type:text" json:"description"`
	PriceMargin float64    `gorm:"default:0" json:"price_margin"`
	Price       int        `json:"price"`
	ResultPrice float64    `json:"result_price"`
	AreComment  int        `gorm:"default:0" json:"are_comment"`
	AreLike     int        `gorm:"default:0" json:"are_like"`
	Comments    *[]Comment `json:"comments,omitempty"`
	Likes       *[]Like    `json:"likes,omitempty"`
	UserID      string     `gorm:"not null" json:"user_id"`
	User        User       `gorm:"foreignkey:UserID;constraint:onUpdate:CASCADE,onDelete:SET NULL" json:"user"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   time.Time  `json:"deleted_at"`
}
