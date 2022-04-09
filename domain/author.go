package domain

import (
	"context"
	"errors"
)

var (
	ErrBookNotFound = errors.New("book not found")
)

type Author struct {
	ID             int    `json:"id" gorm:"primaryKey"`
	Name           string `json:"name" validate:"gte=0,lte=500"`
	Surname        string `json:"surname" validate:"gte=0,lte=500"`
	Email          string `json:"email" gorm:"unique" validate:"email"`
	BooksPublished []Book `json:"books_published" gorm:"many2many:author_books;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" validate:"dive"`
}
type AuthorBooks struct {
	BookID   int `gorm:"primaryKey" column:"book_id"`
	AuthorID int `gorm:"primaryKey" column:"author_id"`
}

type AuthorService interface {
	Create(ctx context.Context, books []string, author *Author) error
	Get(ctx context.Context, id string) (Author, error)
	GetByFilter(ctx context.Context, filter, filterValue string) ([]Author, error)
	Update(ctx context.Context, id string, author *Author, updatedAuthor Author, booksPublished []string) error
	Delete(ctx context.Context, id string, author *Author) error
}

type AuthorRepository interface {
	Create(ctx context.Context, author *Author) error
	Update(ctx context.Context, author *Author, updatedAuthor Author) error
	Get(ctx context.Context, id string) (Author, error)
	Delete(ctx context.Context, id string, author *Author) error
	GetByFilter(ctx context.Context, filter, filterValue string) ([]Author, error)
}

type AuthorBooksRepository interface {
	Create(ctx context.Context, author *Author, authorBooks []Book) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, author *Author, authorBooks []Book) error
}
