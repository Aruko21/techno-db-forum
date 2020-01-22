package postRepository

import (
	"github.com/jackc/pgx"
	"github.com/soulphazed/techno-db-forum/internal/app/post"
	"github.com/soulphazed/techno-db-forum/internal/model"
	"strconv"
)

type PostRepository struct {
	db *pgx.ConnPool
}

func NewPostRepository(db *pgx.ConnPool) post.Repository {
	return &PostRepository{db}
}

func (repo PostRepository) FindById(id string, includeUser, includeForum, includeThread bool) (*model.PostFull, error) {
	postObj := &model.PostFull{}
	postObj.Post = &model.Post{}

	id2, _ := strconv.Atoi(id)

	if err := repo.db.QueryRow(`
		SELECT P.author, P.created, P.forum, P.id, P.message, P.thread, P.isedited, P.parent
			FROM posts AS P
			WHERE P.id = $1
		`,
		id2,
	).Scan(
		&postObj.Post.Author,
		&postObj.Post.Created,
		&postObj.Post.Forum,
		&postObj.Post.ID,
		&postObj.Post.Message,
		&postObj.Post.Thread,
		&postObj.Post.IsEdited,
		&postObj.Post.Parent,
	); err != nil {
		return nil, err
	}

	if includeUser {
		postObj.Author = &model.User{}
		if err := repo.db.QueryRow(`
			SELECT U.about, U.email, U.fullname, U.nickname
				FROM users AS U
				WHERE LOWER(U.nickname) = LOWER($1)
			`,
			postObj.Post.Author,
		).Scan(
			&postObj.Author.About,
			&postObj.Author.Email,
			&postObj.Author.Fullname,
			&postObj.Author.Nickname,
		); err != nil {
			return nil, err
		}
	}

	if includeForum {
		postObj.Forum = &model.Forum{}
		if err := repo.db.QueryRow(`
			SELECT F.author, F.title, F.slug, F.posts, F.threads
				FROM forums AS F
				WHERE LOWER(F.slug) = LOWER($1)
			`,
			postObj.Post.Forum,
		).Scan(
			&postObj.Forum.Author,
			&postObj.Forum.Title,
			&postObj.Forum.Slug,
			&postObj.Forum.Posts,
			&postObj.Forum.Threads,
		); err != nil {
			return nil, err
		}
	}

	if includeThread {
		postObj.Thread = &model.Thread{}
		if err := repo.db.QueryRow(`
			SELECT T.forum, T.slug, T.title, T.author, T.message, T.id, T.created, T.votes
				FROM threads AS T
				WHERE T.id = $1
			`,
			postObj.Post.Thread,
		).Scan(
			&postObj.Thread.Forum,
			&postObj.Thread.Slug,
			&postObj.Thread.Title,
			&postObj.Thread.Author,
			&postObj.Thread.Message,
			&postObj.Thread.ID,
			&postObj.Thread.Created,
			&postObj.Thread.Votes,
		); err != nil {
			return nil, err
		}
	}

	return postObj, nil
}

func (repo PostRepository) Update(id string, message string) (*model.Post, error) {
	postObj := &model.Post{}

	id2, _ := strconv.Atoi(id)

	if err := repo.db.QueryRow(`
		UPDATE posts
			SET message = $2,
				isEdited = TRUE
			WHERE id = $1
		RETURNING author, created, forum, id, message, thread, isEdited
		`,
		id2,
		message,
	).Scan(
		&postObj.Author,
		&postObj.Created,
		&postObj.Forum,
		&postObj.ID,
		&postObj.Message,
		&postObj.Thread,
		&postObj.IsEdited,
	); err != nil {
		return nil, err
	}

	return postObj, nil
}
