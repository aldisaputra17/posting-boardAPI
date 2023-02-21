package dummy

import (
	"github.com/aldisaputra17/posting-board/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateCrypto(db *gorm.DB) {
	// var crypto *entities.Crypto
	id, _ := uuid.NewRandom()

	crypto := &entities.Crypto{
		ID:    id,
		Asset: "BTC",
		Value: 12000,
	}
	db.Create(crypto)
}
