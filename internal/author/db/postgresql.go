package db

import (
	"context"
	"errors"
	"fmt"
	"rest-api-test/internal/apperror"
	"rest-api-test/internal/author"
	"rest-api-test/pkg/client/postgresql"
	"rest-api-test/pkg/logging"
	"strings"

	"github.com/jackc/pgconn"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepository(client postgresql.Client, logger *logging.Logger) author.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func formatQuery(q string) string {
	return fmt.Sprintf("SQL Query: %s", strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", ""))
}

func (r *repository) Create(ctx context.Context, author *author.Author) error {
	q := `
			INSERT INTO author
				(name)
			VALUES
				($1)
			RETURNING id
		`
	r.logger.Trace(formatQuery(q))
	if err := r.client.QueryRow(ctx, q, author.Name).Scan(&author.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf((fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context) (u []author.Author, err error) {
	q := `
		SELECT id, name FROM author;
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	authors := make([]author.Author, 0)

	for rows.Next() {
		var ath author.Author

		err = rows.Scan(&ath.ID, &ath.Name)
		if err != nil {
			return nil, err
		}

		authors = append(authors, ath)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (author.Author, error) {
	q := `
		SELECT id, name FROM author WHERE id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var ath author.Author
	err := r.client.QueryRow(ctx, q, id).Scan(&ath.ID, &ath.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf((fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())))
			r.logger.Error(newErr)
			return author.Author{}, newErr
		}
		return author.Author{}, err
	}

	return ath, nil
}

func (r *repository) Update(ctx context.Context, author author.Author) error {
	q := `
		UPDATE author SET name = $1 where id = $2
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	result, err := r.client.Exec(ctx, q, author.Name, author.ID)
	if err != nil {
		return err
	}
	affected := result.RowsAffected()
	if affected == 0 {
		return apperror.ErrNotFound
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return nil
}
