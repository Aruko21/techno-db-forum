package post

import "github.com/soulphazed/techno-db-forum/internal/model"

type Repository interface {
	FindById(id string, includeUser, includeForum, includeThread bool) (*model.PostFull, error)
	Update(id string, message string) (*model.Post, error)
}
