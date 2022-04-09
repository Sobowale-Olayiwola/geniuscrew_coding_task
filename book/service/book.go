package service

import (
	"context"
	"geniuscrew/domain"
)

type bookService struct {
	bookRepository domain.BookRepository
}

func NewBookService(b domain.BookRepository) domain.BookService {
	return &bookService{bookRepository: b}
}

func (p *bookService) Create(ctx context.Context, book *domain.Book) error {

	err := p.bookRepository.Create(ctx, book)
	return err
}

func (p *bookService) Get(ctx context.Context, id string) (domain.Book, error) {
	book, err := p.bookRepository.Get(ctx, id)
	return book, err
}

func (p *bookService) GetByFilter(ctx context.Context, filter, filterValue string) ([]domain.Book, error) {

	book, err := p.bookRepository.GetByFilter(ctx, filter, filterValue)
	if len(book) == 0 {
		return book, domain.ErrBookNotFound
	}
	return book, err
}

func (p *bookService) Update(ctx context.Context, id string, book *domain.Book, updatedBook domain.Book) error {

	err := p.bookRepository.Update(ctx, book, updatedBook)
	return err
}

func (p *bookService) Delete(ctx context.Context, id string, book *domain.Book) error {
	err := p.bookRepository.Delete(ctx, id, book)
	return err
}
