package result

import "go-starter-template/pkg/domain/entities"

type UserResult struct {
	ID        int
	Email     string
	CreatedAt string
	UpdatedAt string
}

func NewUserResult(user *entities.User) *UserResult {
	return &UserResult{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.ToString(),
		UpdatedAt: user.UpdatedAt.ToString(),
	}
}
