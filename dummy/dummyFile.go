package dummy

import (
	"github.com/aldisaputra17/posting-board/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateCrypto(db *gorm.DB) {
	// var crypto *entities.Crypto
	id, _ := uuid.NewRandom()

	for i := 0; i < 100; i++ {
		crypto := &entities.Crypto{
			ID:    id,
			Asset: "Doge",
			Value: 19000000,
		}
		db.Create(crypto)
	}
}
