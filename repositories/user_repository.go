package repositories

import (
	"database/sql"
	"marketplace-app/entities"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(username string) (*entities.User, error) {
	query := "SELECT id, username, name, password FROM users WHERE username = $1 AND deleted_at IS NULL"
	row := r.db.QueryRow(query, username)

	user := &entities.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Create(user *entities.User) error {
	query := "INSERT INTO users (id, username, name, password) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, user.ID, user.Username, user.Name, user.Password)
	return err
}
