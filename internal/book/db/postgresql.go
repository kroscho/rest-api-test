package db

import (
	"context"
	"fmt"
	"rest-api-test/internal/book"
	"rest-api-test/pkg/client/postgresql"
	"rest-api-test/pkg/logging"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

/*
func NewRepository(client postgresql.Client, logger *logging.Logger) author.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
*/

func formatQuery(q string) string {
	return fmt.Sprintf("SQL Query: %s", strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", ""))
}

func (r *repository) FindAll(ctx context.Context) (u []book.Book, err error) {
	q := `
		SELECT id, name, age FROM book	
	`
	r.logger.Trace("SQL Query: %s", formatQuery(q))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	books := make([]book.Book, 0)

	for rows.Next() {
		var bk Book

		err = rows.Scan(&bk.ID, &bk.Name, &bk.Age)
		if err != nil {
			return nil, err
		}

		/*
			sq := `
				SELECT
					a.id, a.name
				FROM book_authors ba
				JOIN author a on a.id = ba.author_id
				WHERE book_id = $1;
			`

			authorRows, err := r.client.Query(ctx, sq, bk.ID)
			if err != nil {
				return nil, err
			}

			authors := make([]author.Author, 0)

			for authorRows.Next() {
				var ath author.Author

				err = authorRows.Scan(&ath.ID, &ath.Name)
				if err != nil {
					return nil, err
				}
				authors = append(authors, ath)
			}
			bk.Authors = authors
		*/

		books = append(books, bk.ToDomain())
	}

	return books, nil
}
