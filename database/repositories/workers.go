package repositories

import (
	"context"

	"github.com/Rasikrr/core/database"
)

type WorkerModel struct {
	ID      int
	Name    string
	Man     bool
	Balance float64
	Status  string
}

type WorkersRepository interface {
	Create(ctx context.Context, worker WorkerModel) error
	Delete(ctx context.Context, id int) error
}

type workersRepository struct {
	db *database.Postgres
}

func NewWorkersRepository(db *database.Postgres) WorkersRepository {
	return &workersRepository{db: db}
}

func (r *workersRepository) Create(ctx context.Context, worker WorkerModel) error {
	_, err := r.db.Exec(ctx, "INSERT INTO workers (name, man, balance, status) VALUES ($1, $2, $3, $4)", worker.Name, worker.Man, worker.Balance, worker.Status)
	return err
}

func (r *workersRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM workers WHERE id = $1", id)
	return err
}
