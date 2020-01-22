package thread

import "github.com/soulphazed/techno-db-forum/internal/model"

type Repository interface {
	CreateThread(thread *model.NewThread) (*model.Thread, error)
	FindByIdOrSlug(id int, slug string) (*model.Thread, error)
	CreatePosts(thread *model.Thread, posts *model.Posts) (*model.Posts, error)
	UpdateThread(id int, slug string, update *model.ThreadUpdate) (*model.Thread, error)
	GetThreadPosts(thread *model.Thread, limit, desc, since, sort string) (model.Posts, error)
	Vote(thread *model.Thread, vote *model.Vote) (*model.Thread, error)
}