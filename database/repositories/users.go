package repositories

import (
	"context"

	"github.com/Rasikrr/core/database"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type UserModel struct {
	Name     string
	Password string
}

type UsersRepository interface {
	Create(ctx context.Context, user UserModel) error
	UpdatePassword(ctx context.Context, user UserModel) error
	GetAll(ctx context.Context) ([]UserModel, error)
}

type usersRepository struct {
	db *database.Postgres
}

func NewUsersRepository(db *database.Postgres) UsersRepository {
	return &usersRepository{db: db}
}

func (r *usersRepository) Create(ctx context.Context, user UserModel) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users (name, password) VALUES ($1, $2)", user.Name, user.Password)
	return err
}

func (r *usersRepository) UpdatePassword(ctx context.Context, user UserModel) error {
	_, err := r.db.Exec(ctx, "UPDATE users SET password = $1 WHERE name = $2", user.Password, user.Name)
	return err
}

func (r *usersRepository) GetAll(ctx context.Context) ([]UserModel, error) {
	var mm []UserModel
	// Используем GetQuerier для автоматического выбора транзакции или pool
	err := pgxscan.Select(ctx, r.db, &mm, "SELECT name, password FROM users")
	return mm, err
}
