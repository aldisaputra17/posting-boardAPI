package dto

type RequestComment struct {
	Comments string `json:"comment" form:"comment" binding:"required"`
	UserID   string `json:"-" form:"user_id,omitempty"`
	PostID   string `json:"post_id" form:"post_id" binding:"required"`
}
