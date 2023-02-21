package dto

type RequestRegister struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
	CryptoID string `json:"crypto_id" form:"crypto_id" binding:"required"`
}
