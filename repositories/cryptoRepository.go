package repositories

import (
	"github.com/aldisaputra17/posting-board/entities"
	"gorm.io/gorm"
)

type CryptoRepository interface {
	FindByID(id string) (*entities.Crypto, error)
}

type cryptoRepository struct {
	connection *gorm.DB
}

func NewCryptoRepository(db *gorm.DB) CryptoRepository {
	return &cryptoRepository{
		connection: db,
	}
}

func (db *cryptoRepository) FindByID(id string) (*entities.Crypto, error) {
	var crypto *entities.Crypto
	res := db.connection.Where("id = ?", id).Find(&crypto)
	if res.Error != nil {
		return nil, res.Error
	}

	return crypto, nil
}
