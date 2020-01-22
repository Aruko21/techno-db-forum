package userUsecase

import (
	"fmt"
	"github.com/soulphazed/techno-db-forum/internal/app/user"
	"github.com/soulphazed/techno-db-forum/internal/model"
)

type UserUsecase struct {
	rep user.Repository
}

func (u UserUsecase) CreateUser(user *model.User) (model.Users, error) {
	users, err := u.rep.Find(user.Nickname, user.Email)

	if len(users) > 0 || err != nil {
		return users, err
	}

	err = u.rep.Create(user)

	if err != nil {
		fmt.Println("CHECK USECASE ", err)
		return nil, err
	}

	return nil, nil
}

func (u UserUsecase) Find(nickname string) (*model.User, error) {
	return u.rep.FindByNickname(nickname)
}

func (u UserUsecase) Update(user *model.User) (*model.User, error, int) {
	userObj, err := u.rep.FindByNickname(user.Nickname)

	if err != nil || userObj == nil {
		return nil, err, 404
	}

	// Если юзер не трогал эти поля - то просто перепрокинуть их
	if user.Email == "" {
		user.Email = userObj.Email
	}
	if user.Fullname == "" {
		user.Fullname = userObj.Fullname
	}
	if user.About == "" {
		user.About = userObj.About
	}

	userObj, err = u.rep.Update(user)

	if err != nil {
		return nil, err, 409
	}

	return userObj, err, 200
}

func NewForumUsecase(rep user.Repository) user.Usecase {
	return &UserUsecase{
		rep: rep,
	}
}
