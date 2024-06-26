package utils

import (
	"context"
	"time"

	pgx "github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
)

const (
	createdAtTimeStr = "2022-08-05 04:29:15.650318"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	Prepare(context.Context, string, string) (*pgconn.StatementDescription, error)

	Close()
}

var SampleCreatedAtTime, _ = time.Parse(LayoutPostgres, createdAtTimeStr)
