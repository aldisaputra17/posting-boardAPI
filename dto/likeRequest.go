package dto

type RequestLike struct {
	Like   int    `json:"like" form:"like" binding:"required,max=1"`
	UserID string `json:"-" form:"user_id,omitempty"`
	PostID string `json:"post_id" form:"post_id" binding:"required"`
}
