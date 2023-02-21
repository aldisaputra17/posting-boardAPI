package services

import (
	"context"
	"time"

	"github.com/aldisaputra17/posting-board/dto"
	"github.com/aldisaputra17/posting-board/entities"
	"github.com/aldisaputra17/posting-board/repositories"
	"github.com/google/uuid"
)

type LikeService interface {
	Add(ctx context.Context, likeReq *dto.RequestLike) (*entities.Like, error)
	ListUser(ctx context.Context, id string) ([]*entities.Like, error)
}

type likeService struct {
	likeRepository repositories.LikeRepository
	contextTimeout time.Duration
}

func NewLikeService(likeRepo repositories.LikeRepository, time time.Duration) LikeService {
	return &likeService{
		likeRepository: likeRepo,
		contextTimeout: time,
	}
}

func (service *likeService) Add(ctx context.Context, likeReq *dto.RequestLike) (*entities.Like, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	likeCreate := &entities.Like{
		ID:        id,
		Likes:     likeReq.Like,
		UserID:    likeReq.UserID,
		PostID:    likeReq.PostID,
		CreatedAt: time.Now(),
	}
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeout)
	defer cancel()

	res, err := service.likeRepository.Add(ctx, likeCreate)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *likeService) ListUser(ctx context.Context, id string) ([]*entities.Like, error) {
	res, err := service.likeRepository.ListUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, err
}
