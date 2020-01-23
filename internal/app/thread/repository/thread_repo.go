package threadRepository

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/jackc/pgx"
	"github.com/soulphazed/techno-db-forum/internal/app/thread"
	"github.com/soulphazed/techno-db-forum/internal/model"
	"strconv"
	"strings"
	"time"
)

type ThreadRepository struct {
	db *pgx.ConnPool
}

func NewThreadRepository(db *pgx.ConnPool) thread.Repository {
	return &ThreadRepository{db}
}

func (repo ThreadRepository) CreateThread(newThread *model.NewThread) (*model.Thread, error) {
	threadObj := &model.Thread{}
	var row *pgx.Row

	if newThread.Created.IsZero() {
		query := `
			INSERT INTO threads (title, message, forum, author, slug)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING slug, title, message, forum, author, created, votes, id
		`
		row = repo.db.QueryRow(
			query,
			newThread.Title,
			newThread.Message,
			newThread.Forum,
			newThread.Author,
			newThread.Slug,
		)
	} else {
		query := `
			INSERT INTO threads (title, message, forum, author, slug, created)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING slug, title, message, forum, author, created, votes, id
		`
		row = repo.db.QueryRow(
			query,
			newThread.Title,
			newThread.Message,
			newThread.Forum,
			newThread.Author,
			newThread.Slug,
			newThread.Created,
		)
	}

	err := row.Scan(
		&threadObj.Slug,
		&threadObj.Title,
		&threadObj.Message,
		&threadObj.Forum,
		&threadObj.Author,
		&threadObj.Created,
		&threadObj.Votes,
		&threadObj.ID,
	)
	if err != nil {
		return nil, err
	}

	_, err = repo.db.Exec(`
		UPDATE forums
			SET threads = threads + 1
			WHERE lower(slug) = lower($1)
		`,
		threadObj.Forum)
	if err != nil {
		return nil, err
	}

	return threadObj, nil
}

func (repo ThreadRepository) UpdateThread(id int, slug string, threadUpdate *model.ThreadUpdate) (*model.Thread, error) {
	threadObj := &model.Thread{}

	err := repo.db.QueryRow(`
		UPDATE threads
			SET title = $1,
				message = $2
			WHERE id=$3 OR LOWER(slug)=LOWER($4)
			RETURNING slug, title, message, forum, author, created, votes, id
		`,
		threadUpdate.Title,
		threadUpdate.Message,
		id,
		slug,
	).Scan(
		&threadObj.Slug,
		&threadObj.Title,
		&threadObj.Message,
		&threadObj.Forum,
		&threadObj.Author,
		&threadObj.Created,
		&threadObj.Votes,
		&threadObj.ID,
	)

	if err != nil {
		return nil, err
	}

	return threadObj, nil
}

func (repo ThreadRepository) FindByIdOrSlug(id int, slug string) (*model.Thread, error) {
	threadObj := &model.Thread{}

	err := repo.db.QueryRow(`
		SELECT slug, title, message, forum, author, created, votes, id
			FROM threads
			WHERE id=$1 OR (LOWER(slug)=LOWER($2) AND slug <> '')
		`,
		id,
		slug,
	).Scan(
		&threadObj.Slug,
		&threadObj.Title,
		&threadObj.Message,
		&threadObj.Forum,
		&threadObj.Author,
		&threadObj.Created,
		&threadObj.Votes,
		&threadObj.ID,
	)

	if err != nil {
		return nil, err
	}

	return threadObj, nil
}

func (repo ThreadRepository) CreatePosts(thread *model.Thread, posts *model.Posts) (*model.Posts, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	sqlStr := `
		INSERT INTO posts(id, parent, thread, forum, author, created, message, path)
			VALUES
	`
	vals := []interface{}{}
	for _, post := range *posts {
		var author string
		err = repo.db.QueryRow(`
			SELECT nickname
				FROM users
				WHERE LOWER(nickname) = LOWER($1)
			`,
			post.Author,
		).Scan(&author)
		// Если хотя бы одного юзера не существует - откатываемся
		if err != nil || author == "" {
			_ = tx.Rollback()
			return nil, errors.New("404")
		}

		if post.Parent == 0 {
			// Создание массива пути с единственным значением -
			// id создаваемого сообщения
			sqlStr += `
				(nextval('posts_id_seq'::regclass), ?, ?, ?, ?, ?, ?, 
				ARRAY[currval(pg_get_serial_sequence('posts', 'id'))::bigint]),`

			vals = append(vals, post.Parent, thread.ID, thread.Forum, post.Author, now, post.Message)
		} else {
			var parentThreadId int32
			err = repo.db.QueryRow(`
				SELECT thread
					FROM posts
					WHERE id = $1
				`,
				post.Parent,
			).Scan(
				&parentThreadId,
			)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}

			if parentThreadId != thread.ID {
				_ = tx.Rollback()
				return nil, errors.New("Parent post was created in another thread")
			}

			// Конкатенация 2-х массивов
			sqlStr += `
				(nextval('posts_id_seq'::regclass), ?, ?, ?, ?, ?, ?, 
				(SELECT path FROM posts WHERE id = ? AND thread = ?) || 
				currval(pg_get_serial_sequence('posts', 'id'))::bigint),`

			vals = append(vals, post.Parent, thread.ID, thread.Forum, post.Author, now, post.Message, post.Parent, thread.ID)
		}

	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	sqlStr += `
		RETURNING  id, parent, thread, forum, author, created, message, isedited
	`

	sqlStr = ReplaceSQL(sqlStr, "?")
	if len(*posts) > 0 {
		rows, err := tx.Query(sqlStr, vals...)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		i := 0
		for rows.Next() {
			err := rows.Scan(
				&(*posts)[i].ID,
				&(*posts)[i].Parent,
				&(*posts)[i].Thread,
				&(*posts)[i].Forum,
				&(*posts)[i].Author,
				&(*posts)[i].Created,
				&(*posts)[i].Message,
				&(*posts)[i].IsEdited,
			)
			i += 1

			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	_, err = tx.Exec(`
		UPDATE forums
			SET posts = posts + $1
			WHERE lower(slug) = lower($2)
		`,
		len(*posts),
		thread.Forum,
	)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (repo ThreadRepository) GetThreadPosts(thread *model.Thread, limit, desc, since, sort string) (model.Posts, error) {
	posts := make(model.Posts, 0)

	var query string

	conditionSign := ">"
	if desc == "desc" {
		conditionSign = "<"
	}

	query = `
		SELECT id, parent, thread, forum, author, created, message, isedited, path
			FROM posts
			WHERE thread = $1
	`

	switch sort {
	case "tree": {
		if since != "" {
			query += fmt.Sprintf(" AND path %s (SELECT path FROM posts WHERE id = %s)", conditionSign, since)
		}
		query += fmt.Sprintf(" ORDER BY path %s ", desc)
		query += fmt.Sprintf(" LIMIT %s", limit)
	}

	case "parent_tree": {
		query += `
			AND path[1] IN ( 
				SELECT id FROM posts
					WHERE thread = $1
		`
		if since != "" {
			query += fmt.Sprintf(" AND id %s (SELECT path[1] FROM posts WHERE id = %s)", conditionSign, since)
		} else {
			// Данное условие нужно, чтобы не писать отдельный индекс на thread, parent, а использовать thread, id, parent
			query += fmt.Sprintf(" AND id > 0")
		}
		// TODO: Здесь возможен индекс, но я не уверен, что он не нагрузит слишком сильно систему при заполнении данных
		query += ` AND parent = 0 `
		query += fmt.Sprintf(" ORDER BY path[1] %s LIMIT %s )", desc, limit )

		if desc == "desc" {
			query += ` ORDER BY path[1] DESC, path`
		} else {
			query += ` ORDER BY path`
		}
	}

	default: {
		if since != "" {
			query += fmt.Sprintf(" AND id %s %s ", conditionSign, since)
		}
		query += fmt.Sprintf(" ORDER BY created %s, id %s LIMIT %s", desc, desc, limit)
	}

	}

	//case "parent_tree" {
	//	query += `
	//		AND path && (SELECT ARRAY (select id from posts WHERE thread = $1 AND parent = 0
	//	`
	//	if since != "" {
	//		query += fmt.Sprintf(" AND path %s (SELECT path[1:1] FROM posts WHERE id = %s) ", conditionSign, since)
	//	}
	//	query += fmt.Sprintf("ORDER BY path[1] %s, path LIMIT %s)) ", desc, limit)
	//	query += fmt.Sprintf("ORDER BY path[1] %s, path ", desc)
	//}

	rows, err := repo.db.Query(query, thread.ID)
	if err != nil {
		fmt.Println("Sort error!: ", err, "\n for SQL: ", query)
		return nil, err
	}

	for rows.Next() {
		p := model.Post{}
		err := rows.Scan(
			&p.ID,
			&p.Parent,
			&p.Thread,
			&p.Forum,
			&p.Author,
			&p.Created,
			&p.Message,
			&p.IsEdited,
			&p.Path,
		)
		if err != nil {
			fmt.Println("Error while scanning path: ", err)
			return nil, err
		}

		posts = append(posts, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (repo ThreadRepository) Vote(thread *model.Thread, vote *model.Vote) (*model.Thread, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO votes(nickname, voice, thread)
			VALUES ($1, $2, $3)
		ON CONFLICT ON CONSTRAINT votes_pkey
			DO UPDATE SET voice = $2
		`,
		vote.Nickname,
		vote.Voice,
		thread.ID,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.QueryRow(`
		SELECT votes
			FROM threads
			WHERE id = $1
		`,
		thread.ID,
	).Scan(
		&thread.Votes,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

