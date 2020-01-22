package forum

import "github.com/soulphazed/techno-db-forum/internal/model"

type Usecase interface {
	CreateForum(forum *model.Forum) (*model.Forum, int, error)
	Find(slug string) (*model.Forum, error)
	CreateThread(slug string, newThread *model.NewThread) (*model.Thread, int, error)
	GetUsersByForum(forumSlug string, params map[string][]string) (model.Users, int, error)
	GetThreadsByForum(forumSlug string, params map[string][]string) (model.Threads, int, error)
}