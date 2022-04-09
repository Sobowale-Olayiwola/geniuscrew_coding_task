package repository

import (
	"context"
	"errors"
	"fmt"
	"geniuscrew/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type mysqlAuthorRepository struct {
	db *gorm.DB
}

func NewMySqlAuthorRepository(db *gorm.DB) domain.AuthorRepository {
	return &mysqlAuthorRepository{db}
}

func (m *mysqlAuthorRepository) Create(ctx context.Context, author *domain.Author) error {
	err := m.db.WithContext(ctx).Create(author).Error
	return err
}

func (m *mysqlAuthorRepository) Get(ctx context.Context, id string) (domain.Author, error) {
	var author domain.Author
	err := m.db.WithContext(ctx).Preload(clause.Associations).Where("id = ?", id).First(&author).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return domain.Author{}, domain.ErrRecordNotFound
		default:
			return domain.Author{}, err
		}
	}
	return author, nil
}

func (m *mysqlAuthorRepository) Update(ctx context.Context, author *domain.Author, updatedAuthor domain.Author) error {
	err := m.db.WithContext(ctx).Model(author).Updates(updatedAuthor).Error
	return err
}

func (m *mysqlAuthorRepository) Delete(ctx context.Context, id string, author *domain.Author) error {
	err := m.db.WithContext(ctx).Select("Book").Where("id = ?", id).Delete(author).Error
	return err
}

func (m *mysqlAuthorRepository) GetByFilter(ctx context.Context, filter, filterValue string) ([]domain.Author, error) {
	var authors []domain.Author
	filterValue = "%" + filterValue + "%"
	err := m.db.WithContext(ctx).Preload(clause.Associations).Where(fmt.Sprintf("%s LIKE ?", filter), filterValue).Find(&authors).Error
	if err != nil {
		return []domain.Author{}, err
	}
	return authors, nil
}
