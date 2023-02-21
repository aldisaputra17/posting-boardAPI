package repositories

import (
	"context"
	"log"

	"github.com/aldisaputra17/posting-board/entities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) (*entities.User, error)
	FindByID(ctx context.Context, id string) (*entities.User, error)
	VerifyCredential(email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	user.Password = hashAndSalt([]byte(user.Password))
	cRr := NewCryptoRepository(db.connection)
	cryptoItem, err := cRr.FindByID(user.CryptoID)
	if err != nil {
		return nil, err
	}
	user.Asset = cryptoItem.Asset
	user.Value = cryptoItem.Value
	res := db.connection.WithContext(ctx).Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (db *userConnection) VerifyCredential(email string, password string) interface{} {
	var user entities.User
	res := db.connection.Where("email = ?", email).Take(&user)
	if res.Error == nil {
		return user
	}
	return nil
}

func (db *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user entities.User
	return db.connection.Where("email = ?", email).Take(&user)
}

func (db *userConnection) FindByID(ctx context.Context, id string) (*entities.User, error) {
	var user *entities.User
	res := db.connection.WithContext(ctx).Where("id = ?", id).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}
