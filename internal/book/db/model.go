package db

import (
	"database/sql"
	"rest-api-test/internal/author"
	"rest-api-test/internal/book"
)

type Book struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Age     sql.NullInt32   `json:"age`
	Authors []author.Author `json:"author"`
}

func (m *Book) ToDomain() book.Book {

	b := book.Book{
		ID:   m.ID,
		Name: m.Name,
	}
	if m.Age.Valid {
		b.Age = int(m.Age.Int32)
	}

	return b
}
