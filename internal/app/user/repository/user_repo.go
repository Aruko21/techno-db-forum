package userRepository

import (
	"github.com/jackc/pgx"
	"github.com/soulphazed/techno-db-forum/internal/app/user"
	"github.com/soulphazed/techno-db-forum/internal/model"
)

type UserRepository struct {
	db *pgx.ConnPool
}

func (u UserRepository) Create(user *model.User) error {
	return u.db.QueryRow(`
		INSERT INTO users (nickname, email, fullname, about)
			VALUES ($1, $2, $3, $4)
			RETURNING nickname;
		`,
		user.Nickname,
		user.Email,
		user.Fullname,
		user.About,
	).Scan(&user.Nickname)
}

func (u UserRepository) FindByNickname(nickname string) (*model.User, error) {
	userObj := &model.User{}

	if err := u.db.QueryRow(`
		SELECT U.nickname, U.email, U.fullname, U.about
			FROM users AS U
			WHERE LOWER(U.nickname) = LOWER($1);
		`,
		nickname,
	).Scan(
		&userObj.Nickname,
		&userObj.Email,
		&userObj.Fullname,
		&userObj.About,
	); err != nil {
		return nil, err
	}

	return userObj, nil
}

func (u UserRepository) Find(nickname string, email string) (model.Users, error) {
	var users model.Users

	rows, err := u.db.Query(`
		SELECT U.nickname, U.about, U.email, U.fullname
			FROM users AS U
			WHERE LOWER(U.nickname) = LOWER($1) OR LOWER(U.email) = LOWER($2);
		`,
		nickname,
		email,
	)

	if err != nil {
		return nil, err
	}

	// rows.Close() вызывается автоматически в этом цикле.
	// Достаточно лишь проверить ошибки, что делается в rows.Err()
	for rows.Next() {
		obj := model.User{}
		err := rows.Scan(
			&obj.Nickname,
			&obj.About,
			&obj.Email,
			&obj.Fullname,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, obj)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u UserRepository) Update(user *model.User) (*model.User, error) {
	if err := u.db.QueryRow(`
		UPDATE users
			SET about = $1,
				email = $2,
				fullname = $3
			WHERE LOWER(nickname) = LOWER($4)
			RETURNING nickname, about, email, fullname
		`,
		user.About,
		user.Email,
		user.Fullname,
		user.Nickname,
	).Scan(
		&user.Nickname,
		&user.About,
		&user.Email,
		&user.Fullname,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserRepository(db *pgx.ConnPool) user.Repository {
	return &UserRepository{db}
}
