package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aldisaputra17/posting-board/dto"
	"github.com/aldisaputra17/posting-board/entities"
	"github.com/aldisaputra17/posting-board/helpers"
	"github.com/aldisaputra17/posting-board/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostService interface {
	Add(ctx context.Context, reqPost *dto.RequestPost) (*entities.Post, error)
	Deleted(ctx context.Context, post entities.Post) error
	ListPrivate(ctx context.Context, userID string) ([]*entities.Post, error)
	ListPublic(ctx context.Context) ([]*entities.Post, error)
	LikeMost(ctx *gin.Context, paginate *entities.Pagination) (helpers.Response, error)
	CommentMost(ctx *gin.Context, paginate *entities.Pagination) (helpers.Response, error)
	FindByUserID(ctx context.Context, userID string) ([]*entities.Post, error)
	IsAllowedToEdit(ctx context.Context, postID string, userID string) bool
}

type postService struct {
	postRepository repositories.PostRepository
	contextTimeout time.Duration
}

func NewPostService(postRepo repositories.PostRepository, time time.Duration) PostService {
	return &postService{
		postRepository: postRepo,
		contextTimeout: time,
	}
}

func (service *postService) Add(ctx context.Context, reqPost *dto.RequestPost) (*entities.Post, error) {
	// var trade entities.Trade
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	postCreate := &entities.Post{
		ID:          id,
		Asset:       reqPost.Asset,
		Fiat:        reqPost.Fiat,
		PriceMargin: reqPost.PriceMargin,
		Description: reqPost.Description,
		UserID:      reqPost.UserID,
		CreatedAt:   time.Now(),
	}
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeout)
	defer cancel()

	res, err := service.postRepository.Add(ctx, postCreate)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *postService) Deleted(ctx context.Context, post entities.Post) error {
	return service.postRepository.Deleted(ctx, post)
}

func (service *postService) ListPublic(ctx context.Context) ([]*entities.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeout)
	defer cancel()
	res, err := service.postRepository.ListPublic(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *postService) ListPrivate(ctx context.Context, userID string) ([]*entities.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeout)
	defer cancel()
	res, err := service.postRepository.ListPrivate(ctx, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *postService) CommentMost(ctx *gin.Context, paginate *entities.Pagination) (helpers.Response, error) {
	operationResult, totalPages := service.postRepository.CommentMost(ctx, paginate)

	if operationResult.Error != nil {
		return helpers.Response{Success: true, Message: operationResult.Error.Error()}, nil
	}

	var data = operationResult.Result.(*entities.Pagination)

	urlPath := ctx.Request.URL.Path

	searchQueryParams := ""

	for _, search := range paginate.Searchs {
		searchQueryParams += fmt.Sprintf("&%s.%s=%s", search.Column, search.Action, search.Query)
	}

	data.FirstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort_post=%s", urlPath, paginate.Limit, 0, paginate.SortByComment) + searchQueryParams
	data.LastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort_post=%s", urlPath, paginate.Limit, totalPages, paginate.SortByComment) + searchQueryParams

	if data.Page > 0 {
		// set previous page pagination response
		data.PreviousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort_post=%s", urlPath, paginate.Limit, data.Page-1, paginate.SortByComment) + searchQueryParams
	}

	if data.Page < totalPages {
		// set next page pagination response
		data.NextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort_post=%s", urlPath, paginate.Limit, data.Page+1, paginate.SortByComment) + searchQueryParams
	}

	if data.Page > totalPages {
		// reset previous page
		data.PreviousPage = ""
	}

	return helpers.BuildResponse(true, "Ok", data), nil

}

func (service *postService) LikeMost(ctx *gin.Context, paginate *entities.Pagination) (helpers.Response, error) {
	operationResult, totalPages := service.postRepository.LikeMost(ctx, paginate)

	if operationResult.Error != nil {
		return helpers.Response{Success: true, Message: operationResult.Error.Error()}, nil
	}

	var data = operationResult.Result.(*entities.Pagination)

	urlPath := ctx.Request.URL.Path

	searchQueryParams := ""

	for _, search := range paginate.Searchs {
		searchQueryParams += fmt.Sprintf("&%s.%s=%s", search.Column, search.Action, search.Query)
	}

	data.FirstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort_post=%s", urlPath, paginate.Limit, 0, paginate.SortByLiked) + searchQueryParams
	data.LastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort_post=%s", urlPath, paginate.Limit, totalPages, paginate.SortByLiked) + searchQueryParams

	if data.Page > 0 {
		// set previous page pagination response
		data.PreviousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort_post=%s", urlPath, paginate.Limit, data.Page-1, paginate.SortByLiked) + searchQueryParams
	}

	if data.Page < totalPages {
		// set next page pagination response
		data.NextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort_post=%s", urlPath, paginate.Limit, data.Page+1, paginate.SortByLiked) + searchQueryParams
	}

	if data.Page > totalPages {
		// reset previous page
		data.PreviousPage = ""
	}

	return helpers.BuildResponse(true, "Ok", data), nil
}

func (service *postService) FindByUserID(ctx context.Context, userID string) ([]*entities.Post, error) {
	res, err := service.postRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, service.contextTimeout)
	defer cancel()

	return res, nil
}

func (service *postService) IsAllowedToEdit(ctx context.Context, postID string, userID string) bool {
	c, err := service.postRepository.FindById(ctx, postID, userID)
	if err != nil {
		return false
	}
	id := fmt.Sprintf("%v", c.UserID)
	return userID == id
}
