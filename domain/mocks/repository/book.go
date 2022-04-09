package repository

import (
	"context"
	"geniuscrew/domain"

	"github.com/stretchr/testify/mock"
)

type BookRepositoryMock struct {
	mock.Mock
}

func (w *BookRepositoryMock) Create(ctx context.Context, book *domain.Book) error {
	output := w.Mock.Called(ctx, book)
	err := output.Error(0)
	return err
}

func (w *BookRepositoryMock) Get(ctx context.Context, id string) (domain.Book, error) {
	output := w.Mock.Called(ctx, id)
	book := output.Get(0)
	err := output.Error(1)
	return book.(domain.Book), err
}

func (w *BookRepositoryMock) GetByFilter(ctx context.Context, filter, filterValue string) ([]domain.Book, error) {
	output := w.Mock.Called(ctx, filter, filterValue)
	book := output.Get(0)
	err := output.Error(1)
	return book.([]domain.Book), err
}

func (w *BookRepositoryMock) GetByISBN(ctx context.Context, field string, filter []string) ([]domain.Book, error) {
	output := w.Mock.Called(ctx, field, filter)
	book := output.Get(0)
	err := output.Error(1)
	return book.([]domain.Book), err
}

func (w *BookRepositoryMock) Update(ctx context.Context, book *domain.Book, updatedBook domain.Book) error {
	output := w.Mock.Called(ctx, book, updatedBook)
	err := output.Error(0)
	return err
}

func (w *BookRepositoryMock) Delete(ctx context.Context, id string, book *domain.Book) error {
	output := w.Mock.Called(ctx, id, book)
	err := output.Error(0)
	return err
}
