package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrDuplicateRecord = errors.New("duplicate record")
)

type Book struct {
	ID                int       `json:"id" gorm:"primaryKey"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	ISBN              string    `json:"ISBN" gorm:"unique"`
	Authors           []Author  `json:"authors" gorm:"many2many:author_books;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	PublicationDate   string    `json:"publication_date"`
	PublishingCompany string    `json:"publishing_company"`
	CreatedAt         time.Time `json:"created_at" `
	UpdatedAt         time.Time `json:"updated_at"`
}

type BookService interface {
	Create(ctx context.Context, book *Book) error
	Get(ctx context.Context, id string) (Book, error)
	GetByFilter(ctx context.Context, filter, filterValue string) ([]Book, error)
	Update(ctx context.Context, id string, book *Book, updatedBook Book) error
	Delete(ctx context.Context, id string, book *Book) error
}

type BookRepository interface {
	Create(ctx context.Context, book *Book) error
	Update(ctx context.Context, book *Book, updatedBook Book) error
	Get(ctx context.Context, id string) (Book, error)
	GetByFilter(ctx context.Context, filter, filterValue string) ([]Book, error)
	Delete(ctx context.Context, id string, book *Book) error
	GetByISBN(ctx context.Context, field string, filter []string) ([]Book, error)
}
