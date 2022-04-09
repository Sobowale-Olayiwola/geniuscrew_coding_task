package service

import (
	"context"
	"geniuscrew/domain"
	"strings"
)

type authorService struct {
	authorRepository     domain.AuthorRepository
	authorBookRepository domain.AuthorBooksRepository
	bookRepository       domain.BookRepository
}

func NewAuthorService(a domain.AuthorRepository, ab domain.AuthorBooksRepository, b domain.BookRepository) domain.AuthorService {
	return &authorService{authorRepository: a, authorBookRepository: ab, bookRepository: b}
}

func (p *authorService) Create(ctx context.Context, books []string, author *domain.Author) error {
	authorBooks, err := p.bookRepository.GetByISBN(ctx, "ISBN", books)
	if err != nil {
		return err
	}
	if len(authorBooks) == 0 {
		return domain.ErrBookNotFound
	}
	err = p.authorRepository.Create(ctx, author)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return domain.ErrDuplicateRecord
		}
		return err
	}
	err = p.authorBookRepository.Create(ctx, author, authorBooks)
	if err != nil {
		return err
	}
	return err
}

func (p *authorService) Get(ctx context.Context, id string) (domain.Author, error) {
	author, err := p.authorRepository.Get(ctx, id)
	return author, err
}

func (p *authorService) GetByFilter(ctx context.Context, filter, filterValue string) ([]domain.Author, error) {

	author, err := p.authorRepository.GetByFilter(ctx, filter, filterValue)
	return author, err
}

func (p *authorService) Update(ctx context.Context, id string, author *domain.Author, updatedAuthor domain.Author, booksPublished []string) error {
	var authorBooks []domain.Book
	var err error
	if len(booksPublished) > 0 {
		authorBooks, err = p.bookRepository.GetByISBN(ctx, "ISBN", booksPublished)
		if err != nil {
			return err
		}
		if len(authorBooks) == 0 {
			return domain.ErrBookNotFound
		}
	}
	err = p.authorRepository.Update(ctx, author, updatedAuthor)
	if err != nil {
		return err
	}
	err = p.authorBookRepository.Update(ctx, id, author, authorBooks)
	if err != nil {
		return err
	}
	return err
}

func (p *authorService) Delete(ctx context.Context, id string, author *domain.Author) error {
	err := p.authorBookRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	err = p.authorRepository.Delete(ctx, id, author)
	return err
}
