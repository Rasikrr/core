package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/Rasikrr/core/config"
	"github.com/Rasikrr/core/database"
	"github.com/Rasikrr/core/database/repositories"
	"github.com/Rasikrr/core/log"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestTXManager(t *testing.T) {
	ctx := context.Background()
	postgres, err := database.NewPostgres(ctx, config.PostgresConfig{
		DSN:                 "postgres://postgres:rasik1234@localhost:5432/test",
		Required:            true,
		MaxConns:            10,
		MinConns:            1,
		MaxIdleConnIdleTime: 1 * time.Minute,
	})
	require.NoError(t, err)

	txManager := database.NewTXManager(postgres.Pool())

	usersRepository := repositories.NewUsersRepository(postgres)
	workersRepository := repositories.NewWorkersRepository(postgres)

	txOpt := pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	}
	users, err := usersRepository.GetAll(ctx)
	require.NoError(t, err)
	for _, user := range users {
		log.Infof(ctx, "user: %+v", user)
	}
	err = txManager.Transaction(ctx, txOpt, func(ctx context.Context) error {
		if err := usersRepository.Create(ctx, repositories.UserModel{
			Name:     "test1",
			Password: "test1",
		}); err != nil {
			return err
		}
		if err = workersRepository.Create(ctx, repositories.WorkerModel{
			ID:      100,
			Name:    "test1",
			Man:     true,
			Balance: 100,
			Status:  "active",
		}); err != nil {
			return err
		}
		users, err := usersRepository.GetAll(ctx)
		require.NoError(t, err)
		for _, user := range users {
			log.Infof(ctx, "user: %+v", user)
		}
		if err = workersRepository.Delete(ctx, 12); err != nil {
			return err
		}
		return nil
	})
	require.NoError(t, err)
}
