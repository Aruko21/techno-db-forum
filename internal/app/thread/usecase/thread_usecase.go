package threadUsecase

import (
	"github.com/soulphazed/techno-db-forum/internal/app/thread"
	"github.com/soulphazed/techno-db-forum/internal/app/user"
	"github.com/soulphazed/techno-db-forum/internal/model"
	"strconv"
)

type ThreadUsecase struct {
	threadRep thread.Repository
	userRep   user.Repository
}

func NewThreadUsecase(t thread.Repository, u user.Repository) thread.Usecase {
	return &ThreadUsecase{
		threadRep: t,
		userRep:   u,
	}
}

func (t ThreadUsecase) FindByIdOrSlug(threadSlugOrId string) (*model.Thread, error) {
	id, _ := strconv.Atoi(threadSlugOrId)

	threadObj, err := t.threadRep.FindByIdOrSlug(id, threadSlugOrId)
	if threadObj == nil || err != nil {
		return nil, err
	}

	return threadObj, nil
}

func (t ThreadUsecase) Vote(threadSlugOrId string, vote *model.Vote) (*model.Thread, error) {
	threadObj, err := t.FindByIdOrSlug(threadSlugOrId)
	if err != nil {
		return nil, err
	}

	threadObj, err = t.threadRep.Vote(threadObj, vote)
	if err != nil {
		return nil, err
	}

	return threadObj, nil
}

func (t ThreadUsecase) GetThreadPosts(threadSlugOrId string, params map[string][]string) (model.Posts, error) {
	threadObj, err := t.FindByIdOrSlug(threadSlugOrId)
	if err != nil {
		return nil, err
	}

	limit := "100"
	if len(params["limit"]) >= 1 {
		limit = params["limit"][0]
	}

	desc := ""
	if len(params["desc"]) >= 1 && params["desc"][0] == "true" {
		desc = "desc"
	}

	since := ""
	if len(params["since"]) >= 1 {
		since = params["since"][0]
	}

	sort := "flat"
	if len(params["sort"]) >= 1 {
		sort = params["sort"][0]
	}

	posts, err := t.threadRep.GetThreadPosts(threadObj, limit, desc, since, sort)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (t ThreadUsecase) CreatePosts(threadSlugOrId string, posts *model.Posts) (*model.Posts, int, error) {
	threadObj, err := t.FindByIdOrSlug(threadSlugOrId)
	if threadObj == nil || err != nil {
		return nil, 404, err
	}

	posts, err = t.threadRep.CreatePosts(threadObj, posts)
	if err != nil {
		if err.Error() == "404" {
			return nil, 404, err
		}

		return nil, 409, err
	}

	return posts, 201, nil
}

func (t ThreadUsecase) UpdateThread(threadSlugOrId string, threadUpdate *model.ThreadUpdate) (*model.Thread, error) {
	id, _ := strconv.Atoi(threadSlugOrId)

	threadObj, err := t.threadRep.FindByIdOrSlug(id, threadSlugOrId)
	if err != nil {
		return nil, err
	}

	if threadUpdate.Title == "" && threadUpdate.Message == "" {
		return threadObj, nil
	}

	if threadUpdate.Title == "" {
		threadUpdate.Title = threadObj.Title
	}

	if threadUpdate.Message == "" {
		threadUpdate.Message = threadObj.Message
	}

	threadObj, err = t.threadRep.UpdateThread(id, threadSlugOrId, threadUpdate)
	if err != nil {
		return nil, err
	}

	return threadObj, nil
}

func makeTree(posts model.Posts) model.Posts {
	tree := make(model.Posts, 0)
	var parent *model.Post

	for _, p := range posts {
		if len(p.Path) == 1 {
			tree = append(tree, p)
			parent = p
		} else if len(p.Path) > 1 {
			if p.Parent == parent.ID {
				//parent.Childs = append(parent.Childs, p)
				tree = append(tree, p)
				p.ParentPointer = parent
				parent = p
			} else {
				for p.Parent != parent.ID {
					parent = parent.ParentPointer
				}
				//parent.Childs = append(parent.Childs, p)
				tree = append(tree, p)
				p.ParentPointer = parent
				parent = p
			}
		}
	}

	return tree
}
