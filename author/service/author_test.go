package service

import (
	"context"
	"errors"
	"geniuscrew/domain"
	"geniuscrew/domain/mocks/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	as := assert.New(t)
	authorRepo := &repository.AuthorRepositoryMock{}
	bookRepo := &repository.BookRepositoryMock{}
	authorBookRepo := &repository.AuthorBooksRepositoryMock{}
	t.Run("happy path: Successfully creates an author", func(t *testing.T) {
		bookRepo.On("GetByISBN", context.Background(), "ISBN", []string{"978160309028"}).Return([]domain.Book{
			{ID: 1,
				ISBN: "978160309028",
			},
		}, nil).Once()
		authorRepo.On("Create", context.Background(), mock.Anything).Return(nil).Once()
		authorBookRepo.On("Create", context.Background(), mock.Anything, mock.Anything).Return(nil).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		err := service.Create(context.Background(), []string{"978160309028"}, &domain.Author{})
		as.NoError(err)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})
	t.Run("input error: author book provided not found", func(t *testing.T) {
		bookRepo.On("GetByISBN", context.Background(), "ISBN", []string{"978160309028"}).Return([]domain.Book{}, nil).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		err := service.Create(context.Background(), []string{"978160309028"}, &domain.Author{})
		as.Error(err)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})
	t.Run("input error: author exists", func(t *testing.T) {
		bookRepo.On("GetByISBN", context.Background(), "ISBN", []string{"978160309028"}).Return([]domain.Book{
			{ID: 1,
				ISBN: "978160309028",
			},
		}, nil).Once()
		authorRepo.On("Create", context.Background(), mock.Anything).Return(errors.New("Duplicate")).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		err := service.Create(context.Background(), []string{"978160309028"}, &domain.Author{})
		as.Error(err)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})

	t.Run("system error: error fetching book", func(t *testing.T) {
		bookRepo.On("GetByISBN", context.Background(), "ISBN", []string{"978160309028"}).Return([]domain.Book{}, errors.New("Something failed")).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		err := service.Create(context.Background(), []string{"978160309028"}, &domain.Author{})
		as.Error(err)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})

	t.Run("system error: database failed in creating author", func(t *testing.T) {
		bookRepo.On("GetByISBN", context.Background(), "ISBN", []string{"978160309028"}).Return([]domain.Book{
			{ID: 1,
				ISBN: "978160309028",
			},
		}, nil).Once()
		authorRepo.On("Create", context.Background(), mock.Anything).Return(errors.New("something failed")).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		err := service.Create(context.Background(), []string{"978160309028"}, &domain.Author{})
		as.Error(err)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})

	t.Run("system error: database failed in populating record in junction table", func(t *testing.T) {
		bookRepo.On("GetByISBN", context.Background(), "ISBN", []string{"978160309028"}).Return([]domain.Book{
			{ID: 1,
				ISBN: "978160309028",
			},
		}, nil).Once()
		authorRepo.On("Create", context.Background(), mock.Anything).Return(nil).Once()
		authorBookRepo.On("Create", context.Background(), mock.Anything, mock.Anything).Return(errors.New("something failed")).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		err := service.Create(context.Background(), []string{"978160309028"}, &domain.Author{})
		as.Error(err)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})
}
