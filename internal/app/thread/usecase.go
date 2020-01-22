package thread

import "github.com/soulphazed/techno-db-forum/internal/model"

type Usecase interface {
	CreatePosts(threadSlugOrId string, posts *model.Posts) (*model.Posts, int, error)
	FindByIdOrSlug(threadSlugOrId string) (*model.Thread, error)
	UpdateThread(threadSlugOrId string, update *model.ThreadUpdate) (*model.Thread, error)
	GetThreadPosts(threadSlugOrId string, params map[string][]string) (model.Posts, error)
	Vote(threadSlugOrId string, vote *model.Vote) (*model.Thread, error)
}
