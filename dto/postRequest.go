package dto

type RequestPost struct {
	Asset       string  `json:"asset" form:"asset" binding:"required"`
	Fiat        string  `json:"fiat" form:"fiat" binding:"required"`
	Description string  `json:"description" form:"description" binding:"required"`
	Price       int     `json:"price" form:"price" binding:"required"`
	UserID      string  `json:"-" form:"user_id,omitempty"`
	PriceMargin float64 `json:"price_margin" form:"price_margin" binding:"required,max=1"`
}
