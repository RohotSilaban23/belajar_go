package repository

import (
	"belajar_go/model"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(user *model.User) (*model.User, error) {
	query := `INSERT INTO users(name, email, password) VALUES (?, ?, ?)`

	result, err := r.DB.Exec(
		query,
		user.Name,
		user.Email,
		user.Password,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	user.ID = id

	return user, nil
}

func (r *UserRepository) DeleteUser(id int64) error {
	query := `
        DELETE FROM users
        WHERE id = ?
    `

	_, err := r.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByID(
	id int64,
) (*model.User, error) {

	user := &model.User{}
	err := r.DB.QueryRow(
		`
		SELECT id, name, email, password, created_at
		FROM users
		WHERE id = ?
		`,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByEmai(
	email string,
) (*model.User, error) {

	user := &model.User{}
	err := r.DB.QueryRow(
		`
		SELECT * FROM users WHERE email = ?
		`,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) (*model.User, error) {
	query := `
        UPDATE users
        SET name = ?, email = ?, password = ?
        WHERE id = ?
    `

	_, err := r.DB.Exec(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.ID,
	)

	if err != nil {
		return nil, err
	}

	return nil, err
}
