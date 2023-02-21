package repositories

import (
	"context"

	"github.com/aldisaputra17/posting-board/entities"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Add(ctx context.Context, comment *entities.Comment) (*entities.Comment, error)
	List(ctx context.Context) ([]*entities.Comment, error)
	ListUserComment(ctx context.Context, postID string) ([]*entities.Comment, error)
}

type commentConnection struct {
	connection *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentConnection{
		connection: db,
	}
}

func (db *commentConnection) Add(ctx context.Context, comment *entities.Comment) (*entities.Comment, error) {
	pRr := NewPostRepository(db.connection)
	res := db.connection.WithContext(ctx).Create(&comment)
	if res.Error != nil {
		return nil, res.Error
	}
	errUpdate := pRr.UpdateAreComment(ctx, comment.PostID)
	if errUpdate != nil {
		return nil, errUpdate
	}
	return comment, nil
}

func (db *commentConnection) List(ctx context.Context) ([]*entities.Comment, error) {
	var comments []*entities.Comment
	res := db.connection.WithContext(ctx).Preload("Post").Find(&comments)
	if res.Error != nil {
		return nil, res.Error
	}
	return comments, nil
}

func (db *commentConnection) ListUserComment(ctx context.Context, postID string) ([]*entities.Comment, error) {
	var comments []*entities.Comment
	res := db.connection.WithContext(ctx).Preload("User").Find(&comments)
	if res.Error != nil {
		return nil, res.Error
	}
	return comments, nil
}
