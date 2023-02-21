package repositories

import (
	"context"

	"github.com/aldisaputra17/posting-board/entities"
	"gorm.io/gorm"
)

type LikeRepository interface {
	Add(ctx context.Context, like *entities.Like) (*entities.Like, error)
	ListUser(ctx context.Context, postID string) ([]*entities.Like, error)
}

type likeConnection struct {
	connection *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeConnection{
		connection: db,
	}
}

func (db *likeConnection) Add(ctx context.Context, like *entities.Like) (*entities.Like, error) {
	pRr := NewPostRepository(db.connection)
	res := db.connection.WithContext(ctx).Create(&like)
	if res.Error != nil {
		return nil, res.Error
	}
	errUpdate := pRr.UpdateAreLike(ctx, like.PostID)
	if errUpdate != nil {
		return nil, errUpdate
	}
	return like, nil
}

func (db *likeConnection) ListUser(ctx context.Context, postID string) ([]*entities.Like, error) {
	var likes []*entities.Like
	res := db.connection.WithContext(ctx).Where("post_id = ?", postID).Preload("User").Find(&likes)
	if res.Error != nil {
		return nil, res.Error
	}
	return likes, nil
}
