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
