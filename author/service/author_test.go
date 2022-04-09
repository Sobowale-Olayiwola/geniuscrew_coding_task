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

func TestGet(t *testing.T) {
	as := assert.New(t)
	authorRepo := &repository.AuthorRepositoryMock{}
	bookRepo := &repository.BookRepositoryMock{}
	authorBookRepo := &repository.AuthorBooksRepositoryMock{}
	id := "1"
	t.Run("happy path: Successfully fetches an author", func(t *testing.T) {
		authorRepo.On("Get", context.Background(), id).Return(domain.Author{
			Name: "John Doe",
		}, nil).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		author, err := service.Get(context.Background(), id)
		as.NoError(err)
		as.Equal("John Doe", author.Name)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})

	t.Run("input error: author not found", func(t *testing.T) {
		authorRepo.On("Get", context.Background(), id).Return(domain.Author{}, domain.ErrRecordNotFound).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		author, err := service.Get(context.Background(), id)
		as.Error(err)
		as.Equal("", author.Name)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})

	t.Run("system error: Database failed", func(t *testing.T) {
		authorRepo.On("Get", context.Background(), id).Return(domain.Author{}, errors.New("something failed")).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		author, err := service.Get(context.Background(), id)
		as.Error(err)
		as.Equal("", author.Name)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})
}

func TestGetByFilter(t *testing.T) {
	as := assert.New(t)
	authorRepo := &repository.AuthorRepositoryMock{}
	bookRepo := &repository.BookRepositoryMock{}
	authorBookRepo := &repository.AuthorBooksRepositoryMock{}

	t.Run("happy path: Successfully fetches an author by filter", func(t *testing.T) {
		authorRepo.On("GetByFilter", context.Background(), "name", "john").Return([]domain.Author{
			{
				Name: "John Doe",
			},
			{
				Name: "Johnson",
			},
		}, nil).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		authors, err := service.GetByFilter(context.Background(), "name", "john")
		as.NoError(err)
		as.Equal(len(authors), 2)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})

	t.Run("input error: filter doesn't match record", func(t *testing.T) {
		authorRepo.On("GetByFilter", context.Background(), "name", "dgdfhdhj").Return([]domain.Author{}, domain.ErrRecordNotFound).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		authors, err := service.GetByFilter(context.Background(), "name", "dgdfhdhj")
		as.Error(err)
		as.Equal(len(authors), 0)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})

	t.Run("system error: Database failed", func(t *testing.T) {
		authorRepo.On("GetByFilter", context.Background(), "name", "dgdfhdhj").Return([]domain.Author{}, errors.New("something failed")).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		authors, err := service.GetByFilter(context.Background(), "name", "dgdfhdhj")
		as.Error(err)
		as.Equal(len(authors), 0)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	as := assert.New(t)
	id := "1"
	authorRepo := &repository.AuthorRepositoryMock{}
	bookRepo := &repository.BookRepositoryMock{}
	authorBookRepo := &repository.AuthorBooksRepositoryMock{}
	t.Run("happy path: successfully updates an author", func(t *testing.T) {

		bookRepo.On("GetByISBN", context.Background(), "ISBN", []string{"978160309028"}).Return([]domain.Book{
			{
				ISBN:  "978160309028",
				Title: "Testing in golang",
			},
		}, nil).Once()
		authorRepo.On("Update", context.Background(), &domain.Author{}, mock.Anything).Return(nil).Once()
		authorBookRepo.On("Update", context.Background(), id, &domain.Author{}, []domain.Book{
			{
				ISBN:  "978160309028",
				Title: "Testing in golang",
			},
		}).Return(nil).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		err := service.Update(context.Background(), id, &domain.Author{}, domain.Author{}, []string{"978160309028"})
		as.NoError(err)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})

	t.Run("input error: list of books to update doesn't exist", func(t *testing.T) {
		bookRepo.On("GetByISBN", context.Background(), "ISBN", []string{"978160309028"}).Return([]domain.Book{}, nil).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		err := service.Update(context.Background(), id, &domain.Author{}, domain.Author{}, []string{"978160309028"})
		as.Error(err)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})

	t.Run("system error: Database failed in getting list of existing books", func(t *testing.T) {
		bookRepo.On("GetByISBN", context.Background(), "ISBN", []string{"978160309028"}).Return([]domain.Book{}, errors.New("something failed")).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		err := service.Update(context.Background(), id, &domain.Author{}, domain.Author{}, []string{"978160309028"})
		as.Error(err)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})

	t.Run("an error occured while updating author profile", func(t *testing.T) {
		bookRepo.On("GetByISBN", context.Background(), "ISBN", []string{"978160309028"}).Return([]domain.Book{
			{
				ISBN:  "978160309028",
				Title: "Testing in golang",
			},
		}, nil).Once()
		authorRepo.On("Update", context.Background(), &domain.Author{}, mock.Anything).Return(errors.New("an error occured")).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		err := service.Update(context.Background(), id, &domain.Author{}, domain.Author{}, []string{"978160309028"})
		as.Error(err)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})

	t.Run("an error occured while updating junction table for author and books", func(t *testing.T) {
		bookRepo.On("GetByISBN", context.Background(), "ISBN", []string{"978160309028"}).Return([]domain.Book{
			{
				ISBN:  "978160309028",
				Title: "Testing in golang",
			},
		}, nil).Once()
		authorRepo.On("Update", context.Background(), &domain.Author{}, mock.Anything).Return(nil).Once()
		authorBookRepo.On("Update", context.Background(), id, &domain.Author{}, []domain.Book{
			{
				ISBN:  "978160309028",
				Title: "Testing in golang",
			},
		}).Return(errors.New("an error occured")).Once()
		service := NewAuthorService(authorRepo, authorBookRepo, bookRepo)
		err := service.Update(context.Background(), id, &domain.Author{}, domain.Author{}, []string{"978160309028"})
		as.Error(err)
		authorBookRepo.AssertExpectations(t)
		bookRepo.AssertExpectations(t)
		authorRepo.AssertExpectations(t)
	})
}
