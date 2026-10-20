package repository

import (
	"belajar_go/model"
	"database/sql"
	"fmt"
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
	query := `INSERT INTO users(nama, email, password, role, isActive) VALUES (?, ?, ?, ?, ?)`

	result, err := r.DB.Exec(
		query,
		user.Nama,
		user.Email,
		user.Password,
		user.Role,
		false,
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
		SELECT id, nama, email, password, role, isActive
		FROM users
		WHERE id = ?
		`,
		id,
	).Scan(
		&user.ID,
		&user.Nama,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(
	email string,
) (*model.User, error) {

	user := &model.User{}
	err := r.DB.QueryRow(
		`
		SELECT id, nama, email, password, role, isActive FROM users WHERE email = ?
		`,
		email,
	).Scan(
		&user.ID,
		&user.Nama,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) (*model.User, error) {
	query := `
        UPDATE users
        SET nama = ?, email = ?, password = ?, role = ?, isActive = ?
        WHERE id = ?
    `

	_, err := r.DB.Exec(
		query,
		user.Nama,
		user.Email,
		user.Password,
		user.Role,
		user.IsActive,
		user.ID,
	)

	if err != nil {
		return nil, err
	}

	return nil, err
}

func (r *UserRepository) ActivateUser(email string) error {
	result, err := r.DB.Exec(
		`UPDATE users SET isActive = TRUE
         WHERE email = ? AND isActive = FALSE`,
		email,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("akun tidak ditemukan atau sudah aktif")
	}

	return nil
}
