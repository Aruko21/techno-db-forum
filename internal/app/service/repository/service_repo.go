package serviceRepository

import (
	"github.com/jackc/pgx"
	"github.com/soulphazed/techno-db-forum/internal/app/service"
	"github.com/soulphazed/techno-db-forum/internal/model"
)

type ServiceRepository struct {
	db *pgx.ConnPool
}

func NewServiceRepository(db *pgx.ConnPool) service.Repository {
	return &ServiceRepository{db}
}

func (s ServiceRepository) ClearAll() error {
	if _, err := s.db.Exec(`
		TRUNCATE votes, users, posts, threads, forums, forum_users
			RESTART IDENTITY CASCADE
	`); err != nil {
		return err
	}

	return nil
}

func (s ServiceRepository) GetStatus() (*model.Status, error) {
	status := &model.Status{}

	if err := s.db.QueryRow(`
		SELECT (
			SELECT count(*) from forums
		) AS forum, (
			SELECT count(*) from posts
		) AS post, (
			SELECT count(*) from threads
		) AS thread, (
			SELECT count(*) from users
		) AS user
	`,
	).Scan(
		&status.Forum,
		&status.Post,
		&status.Thread,
		&status.User,
	); err != nil {
		return nil, err
	}

	return status, nil
}
