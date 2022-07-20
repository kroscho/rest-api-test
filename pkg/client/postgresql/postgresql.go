package postgresql

import (
	"context"
	"fmt"
	"log"
	"rest-api-test/internal/config"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
}

func NewClient(ctx context.Context, sc config.StorageConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)

	pool, err = pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Fatal("error do with tries postgresql")
	}
	/*
		err = repeatable.DoWithTries(func() error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			pool, err = pgxpool.Connect(ctx, dsn)
			if err != nil {
				return err
			}

			return nil
		}, sc.MaxAttempts, 5*time.Second)

		if err != nil {
			log.Fatal("error do with tries postgresql")
		}
	*/

	return pool, nil
}
