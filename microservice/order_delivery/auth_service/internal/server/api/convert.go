package api

import (
	"auth-service/internal/db"
)

func newUserResponse(user db.USER) userResponse {
	return userResponse{
		Username:          user.USER_NM,
		Email:             user.EMAIL,
		PasswordChangedAt: user.CHG_DT,
		CreatedAt:         user.CREATE_DT,
	}
}
