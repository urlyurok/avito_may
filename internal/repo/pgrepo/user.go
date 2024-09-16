package pgrepo

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo/repoerrs"
	"avitoTech/internal/storage/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	log "log/slog"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) GetByName(ctx context.Context, username string) (entity.User, error) {
	const fn = "repg.pgrepo.user.GetByName"

	sql := `
	SELECT *
	FROM employee
	WHERE username=$1
	`

	var u entity.User
	err := r.Pool.QueryRow(ctx, sql, username).Scan(
		&u.Id,
		&u.Username,
		&u.FirstName,
		&u.LastName,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		log.Error("Error", "fn:", fn, "err:", err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerrs.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("%s: %v", fn, err)
	}

	return u, nil
}

func (r *UserRepo) GetById(ctx context.Context, userId string) (entity.User, error) {
	const fn = "repg.pgrepo.user.GetByName"

	sql := `
	SELECT *
	FROM employee
	WHERE id=$1
	`

	var u entity.User
	err := r.Pool.QueryRow(ctx, sql, userId).Scan(
		&u.Id,
		&u.Username,
		&u.FirstName,
		&u.LastName,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		log.Error("Error", "fn:", fn, "err:", err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerrs.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("%s: %v", fn, err)
	}

	return u, nil
}
