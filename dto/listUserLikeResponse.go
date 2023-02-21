package dto

import (
	"github.com/aldisaputra17/posting-board/entities"
	"github.com/google/uuid"
)

type UserLikeResponse struct {
	ID    uuid.UUID      `json:"id"`
	Users *entities.User `json:"users"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

func NewUserResponse(user *entities.User) []UserResponse {
	return []UserResponse{
		{
			ID:    user.ID,
			Email: user.Email,
		},
	}
}

// func NewUserLikeResponse(like *entities.Like) []UserLikeResponse {
// 	return []UserLikeResponse{
// 		{
// 			ID:    like.ID,
// 			Users: NewUserResponse(like.User[0]),
// 		},
// 	}
// }
