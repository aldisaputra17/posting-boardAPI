package services

import (
	"context"
	"time"

	"github.com/aldisaputra17/posting-board/dto"
	"github.com/aldisaputra17/posting-board/entities"
	"github.com/aldisaputra17/posting-board/repositories"
	"github.com/google/uuid"
)

type CommentService interface {
	Add(ctx context.Context, commentReq *dto.RequestComment) (*entities.Comment, error)
	List(ctx context.Context) ([]*entities.Comment, error)
	ListUserComment(ctx context.Context, postID string) ([]*entities.Comment, error)
}

type commentService struct {
	commentRepository repositories.CommentRepository
	contextTimeout    time.Duration
}

func NewCommentService(commentRepo repositories.CommentRepository, time time.Duration) CommentService {
	return &commentService{
		commentRepository: commentRepo,
		contextTimeout:    time,
	}
}

func (service *commentService) Add(ctx context.Context, commentReq *dto.RequestComment) (*entities.Comment, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	commentCreate := &entities.Comment{
		ID:        id,
		Comments:  commentReq.Comments,
		UserID:    commentReq.UserID,
		PostID:    commentReq.PostID,
		CreatedAt: time.Now(),
	}
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeout)
	defer cancel()

	res, err := service.commentRepository.Add(ctx, commentCreate)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *commentService) List(ctx context.Context) ([]*entities.Comment, error) {
	res, err := service.commentRepository.List(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *commentService) ListUserComment(ctx context.Context, postID string) ([]*entities.Comment, error) {
	res, err := service.commentRepository.ListUserComment(ctx, postID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
