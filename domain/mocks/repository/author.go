package repository

import (
	"context"
	"geniuscrew/domain"

	"github.com/stretchr/testify/mock"
)

type AuthorRepositoryMock struct {
	mock.Mock
}

type AuthorBooksRepositoryMock struct {
	mock.Mock
}

func (w *AuthorBooksRepositoryMock) Create(ctx context.Context, author *domain.Author, authorBooks []domain.Book) error {
	output := w.Mock.Called(ctx, author, authorBooks)
	err := output.Error(0)
	return err
}

func (w *AuthorBooksRepositoryMock) Update(ctx context.Context, id string, author *domain.Author, authorBooks []domain.Book) error {
	output := w.Mock.Called(ctx, id, author, authorBooks)
	err := output.Error(0)
	return err
}

func (w *AuthorBooksRepositoryMock) Delete(ctx context.Context, id string) error {
	output := w.Mock.Called(ctx, id)
	err := output.Error(0)
	return err
}

func (w *AuthorRepositoryMock) Create(ctx context.Context, author *domain.Author) error {
	output := w.Mock.Called(ctx, author)
	err := output.Error(0)
	return err
}

func (w *AuthorRepositoryMock) Get(ctx context.Context, id string) (domain.Author, error) {
	output := w.Mock.Called(ctx, id)
	author := output.Get(0)
	err := output.Error(1)
	return author.(domain.Author), err
}

func (w *AuthorRepositoryMock) GetByFilter(ctx context.Context, filter, filterValue string) ([]domain.Author, error) {
	output := w.Mock.Called(ctx, filter, filterValue)
	author := output.Get(0)
	err := output.Error(1)
	return author.([]domain.Author), err
}

func (w *AuthorRepositoryMock) Update(ctx context.Context, author *domain.Author, updatedAuthor domain.Author) error {
	output := w.Mock.Called(ctx, author, updatedAuthor)
	err := output.Error(0)
	return err
}

func (w *AuthorRepositoryMock) Delete(ctx context.Context, id string, author *domain.Author) error {
	output := w.Mock.Called(ctx, id, author)
	err := output.Error(0)
	return err
}
