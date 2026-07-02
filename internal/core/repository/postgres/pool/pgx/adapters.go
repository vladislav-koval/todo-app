package core_pgx_pool

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	core_postgres_pool "github.com/vladislav-koval/todo-app/internal/core/repository/postgres/pool"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)

	if err != nil {
		return mapErrors(err)
	}

	return nil
}

type commandTag struct {
	pgconn.CommandTag
}

func mapErrors(err error) error {
	const (
		pgxViolatesForeignKeyErrCode = "23503"
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == pgxViolatesForeignKeyErrCode {
		return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrViolatesForeignKey)
	}

	return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrUnknown)
}
