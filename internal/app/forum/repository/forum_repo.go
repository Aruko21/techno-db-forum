package forumRepository

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/soulphazed/techno-db-forum/internal/app/forum"
	"github.com/soulphazed/techno-db-forum/internal/model"
)

type ForumRepository struct {
	db *pgx.ConnPool
}

func (repo ForumRepository) FindForumThreads(forumSlug string, params map[string][]string) (model.Threads, error) {
	limit := "100"

	if len(params["limit"]) >= 1 {
		limit = params["limit"][0]
	}

	desc := ""
	conditionSign := ">="
	if len(params["desc"]) >= 1 && params["desc"][0] == "true" {
		desc = "desc"
		conditionSign = "<="
	}
	since := ""
	if len(params["since"]) >= 1 {
		since = params["since"][0]
	}

	threads := model.Threads{}

	query := `
		SELECT T.id, T.forum, T.author, T.slug, T.created, T.title, T.message, T.votes
			FROM threads AS T
			WHERE LOWER(T.forum) = LOWER($1)
	`

	if since != "" {
		query += fmt.Sprintf(" AND created %s '%s' ", conditionSign, since)
	}

	query += fmt.Sprintf(" ORDER BY created %s LIMIT %s", desc, limit)

	rows, err := repo.db.Query(query, forumSlug)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := model.Thread{}
		err := rows.Scan(
			&t.ID,
			&t.Forum,
			&t.Author,
			&t.Slug,
			&t.Created,
			&t.Title,
			&t.Message,
			&t.Votes,
		)

		if err != nil {
			return nil, err
		}

		threads = append(threads, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return threads, nil
}

func (repo ForumRepository) FindForumUsers(forum *model.Forum, params map[string][]string) (model.Users, error) {
	limit := "100"

	if len(params["limit"]) >= 1 {
		limit = params["limit"][0]
	}

	sinceConditionSign := ">"
	desc := ""
	if len(params["desc"]) >= 1 && params["desc"][0] == "true" {
		desc = "desc"
		sinceConditionSign = "<"
	}

	since := ""
	if len(params["since"]) >= 1 {
		since = params["since"][0]
	}

	users := model.Users{}

	query := `
		SELECT U.nickname, U.email, U.fullname, U.about
			FROM users AS U
			WHERE U.id IN (
				SELECT FU.user_id
					FROM forum_users AS FU
					WHERE FU.forum_id = $1
			)
	`

	if since != "" {
		query += fmt.Sprintf(" AND LOWER(nickname) %s LOWER('%s') ", sinceConditionSign, since)
	}
	query += fmt.Sprintf(" ORDER BY LOWER(nickname) %s LIMIT %s ", desc, limit)

	rows, err := repo.db.Query(query, forum.ID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := model.User{}
		err := rows.Scan(
			&u.Nickname,
			&u.Email,
			&u.Fullname,
			&u.About,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo ForumRepository) Create(forum *model.Forum) error {
	return repo.db.QueryRow(`
		INSERT INTO forums (slug, title, author)
			VALUES ($1, $2, $3)
		RETURNING slug, title, author, posts, threads
		`,
		forum.Slug,
		forum.Title,
		forum.Author,
	).Scan(
		&forum.Slug,
		&forum.Title,
		&forum.Author,
		&forum.Posts,
		&forum.Threads,
	)
}

func (repo ForumRepository) Find(slug string) (*model.Forum, error) {
	forumObj := &model.Forum{}

	if err := repo.db.QueryRow(`
		SELECT F.id, F.slug, F.title, F.author, F.posts, F.threads
			FROM forums AS F
			WHERE LOWER(slug) = LOWER($1)
		`,
		slug,
	).Scan(
		&forumObj.ID,
		&forumObj.Slug,
		&forumObj.Title,
		&forumObj.Author,
		&forumObj.Posts,
		&forumObj.Threads,
	); err != nil {
		return nil, err
	}

	return forumObj, nil
}

func NewForumRepository(db *pgx.ConnPool) forum.Repository {
	return &ForumRepository{db}
}
