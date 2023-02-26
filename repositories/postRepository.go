package repositories

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/aldisaputra17/posting-board/entities"
	"github.com/aldisaputra17/posting-board/helpers"
	"gorm.io/gorm"
)

type PostRepository interface {
	Add(ctx context.Context, post *entities.Post) (*entities.Post, error)
	Deleted(ctx context.Context, post entities.Post) error
	FindById(ctx context.Context, id string, userID string) (*entities.Post, error)
	FindByUserID(ctx context.Context, userID string) ([]*entities.Post, error)
	CommentMost(ctx context.Context, paginate *entities.Pagination) (helpers.PaginationResult, int)
	LikeMost(ctx context.Context, paginate *entities.Pagination) (helpers.PaginationResult, int)
	ListPrivate(ctx context.Context, userID string) ([]*entities.Post, error)
	ListPublic(ctx context.Context) ([]*entities.Post, error)
	UpdateAreComment(ctx context.Context, id string) error
	UpdateAreLike(ctx context.Context, id string) error
}

type postConnection struct {
	connection *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postConnection{
		connection: db,
	}
}

func (db *postConnection) Add(ctx context.Context, post *entities.Post) (*entities.Post, error) {
	uRr := NewUserRepository(db.connection)
	userItem, err := uRr.FindByID(ctx, post.UserID)
	if err != nil {
		return nil, err
	}
	post.ResultPrice = float64(post.Price) * post.PriceMargin
	res := db.connection.WithContext(ctx).Preload("User").Preload("Comments").Preload("Likes").Create(&post)
	if res.Error != nil {
		return nil, res.Error
	}
	db.connection.Model(&userItem).Update("value", gorm.Expr("value - ?", post.Price))
	return post, nil
}

func (db *postConnection) Deleted(ctx context.Context, post entities.Post) error {
	res := db.connection.WithContext(ctx).Delete(post)
	if res.Error != nil {
		return nil
	}
	return nil
}

func (db *postConnection) ListPublic(ctx context.Context) ([]*entities.Post, error) {
	var posts []*entities.Post
	res := db.connection.WithContext(ctx).Preload("User").Preload("Comments").Preload("Likes").Find(&posts)
	if res.Error != nil {
		return nil, res.Error
	}
	return posts, nil
}

func (db *postConnection) ListPrivate(ctx context.Context, userID string) ([]*entities.Post, error) {
	var posts []*entities.Post
	res := db.connection.WithContext(ctx).Where("user_id = ?", userID).Preload("User").Preload("Comments.User").Preload("Likes").Find(&posts)
	if res.Error != nil {
		return nil, res.Error
	}
	return posts, nil
}

func (db *postConnection) CommentMost(ctx context.Context, paginate *entities.Pagination) (helpers.PaginationResult, int) {
	var (
		post       []*entities.Post
		totalRows  int64
		totalPages int
		fromRow    int
		toRow      int
	)

	offset := paginate.Page * paginate.Limit

	find := db.connection.Limit(paginate.Limit).Offset(offset).Order(paginate.SortByComment)
	searchs := paginate.Searchs

	for _, value := range searchs {
		column := value.Column
		action := value.Action
		query := value.Query

		switch action {
		case "equals":
			whereQuery := fmt.Sprintf("%s = ?", column)
			find = find.Where(whereQuery, query)
		case "contains":
			whereQuery := fmt.Sprintf("%s LIKE ?", column)
			find = find.Where(whereQuery, "%"+query+"%")
		case "in":
			whereQuery := fmt.Sprintf("%s IN (?)", column)
			queryArray := strings.Split(query, ",")
			find = find.Where(whereQuery, queryArray)
		}
	}
	find = find.Preload("User").Preload("Comments").Preload("Likes").Find(&post)
	errFind := find.Error

	if errFind != nil {
		return helpers.PaginationResult{Error: errFind}, totalPages
	}

	paginate.Rows = post

	errCount := db.connection.Model(&entities.Post{}).Count(&totalRows).Error

	if errCount != nil {
		return helpers.PaginationResult{Error: errCount}, totalPages
	}

	paginate.TotalRows = int(totalRows)

	totalPages = int(math.Ceil(float64(totalRows)/float64(paginate.Limit))) - 1

	if paginate.Page == 0 {
		fromRow = 1
		toRow = paginate.Limit
	} else {
		if paginate.Page <= totalPages {
			fromRow = paginate.Page*paginate.Limit + 1
			toRow = (paginate.Page + 1) * paginate.Limit
		}
	}

	if toRow > int(totalRows) {
		toRow = int(totalRows)
	}

	paginate.FromRow = fromRow
	paginate.ToRow = toRow

	return helpers.PaginationResult{Result: paginate}, totalPages
}

func (db *postConnection) LikeMost(ctx context.Context, paginate *entities.Pagination) (helpers.PaginationResult, int) {
	var (
		post       []*entities.Post
		totalRows  int64
		totalPages int
		fromRow    int
		toRow      int
	)

	offset := paginate.Page * paginate.Limit

	find := db.connection.Limit(paginate.Limit).Offset(offset).Order(paginate.SortByLiked)
	searchs := paginate.Searchs

	for _, value := range searchs {
		column := value.Column
		action := value.Action
		query := value.Query

		switch action {
		case "equals":
			whereQuery := fmt.Sprintf("%s = ?", column)
			find = find.Where(whereQuery, query)
		case "contains":
			whereQuery := fmt.Sprintf("%s LIKE ?", column)
			find = find.Where(whereQuery, "%"+query+"%")
		case "in":
			whereQuery := fmt.Sprintf("%s IN (?)", column)
			queryArray := strings.Split(query, ",")
			find = find.Where(whereQuery, queryArray)
		}
	}
	find = find.Preload("User").Preload("Comments").Preload("Likes").Find(&post)
	errFind := find.Error

	if errFind != nil {
		return helpers.PaginationResult{Error: errFind}, totalPages
	}

	paginate.Rows = post

	errCount := db.connection.Model(&entities.Post{}).Count(&totalRows).Error

	if errCount != nil {
		return helpers.PaginationResult{Error: errCount}, totalPages
	}

	paginate.TotalRows = int(totalRows)

	totalPages = int(math.Ceil(float64(totalRows)/float64(paginate.Limit))) - 1

	if paginate.Page == 0 {
		fromRow = 1
		toRow = paginate.Limit
	} else {
		if paginate.Page <= totalPages {
			fromRow = paginate.Page*paginate.Limit + 1
			toRow = (paginate.Page + 1) * paginate.Limit
		}
	}

	if toRow > int(totalRows) {
		toRow = int(totalRows)
	}

	paginate.FromRow = fromRow
	paginate.ToRow = toRow

	return helpers.PaginationResult{Result: paginate}, totalPages
}

func (db *postConnection) FindByUserID(ctx context.Context, userID string) ([]*entities.Post, error) {
	var posts []*entities.Post
	res := db.connection.WithContext(ctx).Where("user_id = ?", userID).Find(&posts)
	if res.Error != nil {
		return nil, res.Error
	}
	return posts, nil
}

func (db *postConnection) FindById(ctx context.Context, id string, userID string) (*entities.Post, error) {
	var post *entities.Post
	res := db.connection.WithContext(ctx).Where("id = ? and user_id = ?", id, userID).Find(&post)
	if res.Error != nil {
		return nil, res.Error
	}
	return post, nil
}

func (db *postConnection) UpdateAreComment(ctx context.Context, id string) error {
	var post *entities.Post
	res := db.connection.WithContext(ctx).Model(&post).Where("id = ?", id).Update("are_comment", gorm.Expr("are_comment + ?", 1))
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (db *postConnection) UpdateAreLike(ctx context.Context, id string) error {
	var post *entities.Post
	res := db.connection.WithContext(ctx).Model(&post).Where("id = ?", id).Update("are_like", gorm.Expr("are_like + ?", 1))
	if res.Error != nil {
		return res.Error
	}
	return nil
}
