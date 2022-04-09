package repository

import (
	"context"
	"errors"
	"fmt"
	"geniuscrew/domain"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type mysqlBookRepository struct {
	db *gorm.DB
}

func NewMySqlBookRepository(db *gorm.DB) domain.BookRepository {
	return &mysqlBookRepository{db}
}

func (m *mysqlBookRepository) Create(ctx context.Context, book *domain.Book) error {
	err := m.db.WithContext(ctx).Create(book).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return domain.ErrDuplicateRecord
		}
		return err
	}
	return nil
}

func (m *mysqlBookRepository) Get(ctx context.Context, id string) (domain.Book, error) {
	var book domain.Book
	err := m.db.WithContext(ctx).Preload(clause.Associations).Where("id = ?", id).Find(&book).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return domain.Book{}, domain.ErrRecordNotFound
		default:
			return domain.Book{}, err
		}
	}
	return book, nil
}

func (m *mysqlBookRepository) Update(ctx context.Context, book *domain.Book, updatedBook domain.Book) error {
	err := m.db.WithContext(ctx).Model(book).Updates(updatedBook).Error
	return err
}

func (m *mysqlBookRepository) Delete(ctx context.Context, id string, book *domain.Book) error {
	err := m.db.WithContext(ctx).Where("id = ?", id).Delete(book).Error
	return err
}

func (m *mysqlBookRepository) GetByFilter(ctx context.Context, filter, filterValue string) ([]domain.Book, error) {
	var books []domain.Book
	filterValue = "%" + filterValue + "%"
	err := m.db.WithContext(ctx).Preload(clause.Associations).Where(fmt.Sprintf("%s LIKE ?", filter), filterValue).Find(&books).Error
	if err != nil {
		return []domain.Book{}, err
	}
	return books, nil
}

func (m *mysqlBookRepository) GetByISBN(ctx context.Context, field string, filter []string) ([]domain.Book, error) {
	var books []domain.Book
	err := m.db.WithContext(ctx).Where(fmt.Sprintf("%s IN ?", field), filter).Find(&books).Error
	if err != nil {
		return []domain.Book{}, err
	}
	return books, nil
}
