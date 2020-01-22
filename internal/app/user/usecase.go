package user

import "github.com/soulphazed/techno-db-forum/internal/model"

type Usecase interface {
	CreateUser(user *model.User) (model.Users, error)
	Find(nickname string) (*model.User, error)
	Update(user *model.User) (*model.User, error, int)
}