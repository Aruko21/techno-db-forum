package user

import "github.com/soulphazed/techno-db-forum/internal/model"

type Repository interface {
	Create(forum *model.User) error
	FindByNickname(nickname string) (*model.User, error)
	Find(nickname string, email string) (model.Users, error)
	Update(user *model.User) (*model.User, error)
}