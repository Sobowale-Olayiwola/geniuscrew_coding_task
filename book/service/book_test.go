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
	bookRepo := &repository.BookRepositoryMock{}
	t.Run("happy path: Successfully creates a book", func(t *testing.T) {
		bookRepo.On("Create", context.Background(), mock.Anything).Return(nil).Once()
		service := NewBookService(bookRepo)
		err := service.Create(context.Background(), &domain.Book{})
		as.NoError(err)
		bookRepo.AssertExpectations(t)
	})

	t.Run("input error: Duplicate book", func(t *testing.T) {
		bookRepo.On("Create", context.Background(), mock.Anything).Return(domain.ErrDuplicateRecord).Once()
		service := NewBookService(bookRepo)
		err := service.Create(context.Background(), &domain.Book{})
		as.Error(err)
		bookRepo.AssertExpectations(t)
	})

	t.Run("system error: Database failed", func(t *testing.T) {
		bookRepo.On("Create", context.Background(), mock.Anything).Return(errors.New("Internal error")).Once()
		service := NewBookService(bookRepo)
		err := service.Create(context.Background(), &domain.Book{})
		as.Error(err)
		bookRepo.AssertExpectations(t)
	})
}

func TestGet(t *testing.T) {
	as := assert.New(t)
	bookRepo := &repository.BookRepositoryMock{}
	id := "1"
	t.Run("happy path: Successfully fetches a book", func(t *testing.T) {
		bookRepo.On("Get", context.Background(), id).Return(domain.Book{
			ID:          1,
			ISBN:        "978160309028",
			Title:       "About test",
			Description: "How to write unit test",
		}, nil).Once()
		service := NewBookService(bookRepo)
		book, err := service.Get(context.Background(), id)
		as.NoError(err)
		as.Equal("978160309028", book.ISBN)
		bookRepo.AssertExpectations(t)
	})

	t.Run("input error: Book not found", func(t *testing.T) {
		bookRepo.On("Get", context.Background(), id).Return(domain.Book{}, domain.ErrRecordNotFound).Once()
		service := NewBookService(bookRepo)
		book, err := service.Get(context.Background(), id)
		as.Error(err)
		as.Equal("", book.ISBN)
		bookRepo.AssertExpectations(t)
	})

	t.Run("system error: Database failed", func(t *testing.T) {
		bookRepo.On("Get", context.Background(), id).Return(domain.Book{}, errors.New("Something failed")).Once()
		service := NewBookService(bookRepo)
		book, err := service.Get(context.Background(), id)
		as.Error(err)
		as.Equal("", book.ISBN)
		bookRepo.AssertExpectations(t)
	})
}

func TestGetByFilter(t *testing.T) {
	as := assert.New(t)
	bookRepo := &repository.BookRepositoryMock{}
	filter, filterValue := "title", "test"
	t.Run("happy path: Successfully fetches book by filter", func(t *testing.T) {
		bookRepo.On("GetByFilter", context.Background(), filter, filterValue).Return([]domain.Book{
			{
				ID:          1,
				ISBN:        "978160309028",
				Title:       "About test",
				Description: "How to write unit test",
			},
			{
				ID:          2,
				ISBN:        "978160309029",
				Title:       "Testing in golang",
				Description: "Creating a go file with the _test in filename",
			},
		}, nil).Once()
		service := NewBookService(bookRepo)
		book, err := service.GetByFilter(context.Background(), filter, filterValue)
		as.NoError(err)
		as.Equal(len(book), 2)
		bookRepo.AssertExpectations(t)
	})

	t.Run("input error: Book matches not found", func(t *testing.T) {
		bookRepo.On("GetByFilter", context.Background(), filter, filterValue).Return([]domain.Book{}, domain.ErrBookNotFound).Once()
		service := NewBookService(bookRepo)
		book, err := service.GetByFilter(context.Background(), filter, filterValue)
		as.Error(err)
		as.Equal(len(book), 0)
		bookRepo.AssertExpectations(t)
	})

	t.Run("system error: Database failed", func(t *testing.T) {
		bookRepo.On("GetByFilter", context.Background(), filter, filterValue).Return([]domain.Book{}, errors.New("Something failed")).Once()
		service := NewBookService(bookRepo)
		book, err := service.GetByFilter(context.Background(), filter, filterValue)
		as.Error(err)
		as.Equal(len(book), 0)
		bookRepo.AssertExpectations(t)
	})
}
