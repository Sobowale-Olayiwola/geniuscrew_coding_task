package repository

import (
	"context"
	"geniuscrew/domain"

	"gorm.io/gorm"
)

type mysqlAuthorBooksRepository struct {
	db *gorm.DB
}

func NewMySqlAuthorBooksRepository(db *gorm.DB) domain.AuthorBooksRepository {
	return &mysqlAuthorBooksRepository{db}
}

func (m *mysqlAuthorBooksRepository) Create(ctx context.Context, author *domain.Author, authorBooks []domain.Book) error {
	var err error
	for _, book := range authorBooks {
		err = m.db.WithContext(ctx).Create(&domain.AuthorBooks{
			BookID:   book.ID,
			AuthorID: author.ID,
		}).Error
		if err != nil {
			m.db.Rollback()
			return err
		}
	}
	return err
}

func (m *mysqlAuthorBooksRepository) Update(ctx context.Context, id string, author *domain.Author, authorBooks []domain.Book) error {
	var err error
	err = m.db.WithContext(ctx).Where("author_id = ?", id).Delete(&domain.AuthorBooks{}).Error
	if err != nil {
		return err
	}
	for _, book := range authorBooks {
		err = m.db.WithContext(ctx).Create(&domain.AuthorBooks{
			BookID:   book.ID,
			AuthorID: author.ID,
		}).Error
		if err != nil {
			m.db.Rollback()
			return err
		}
	}
	return err
}

func (m *mysqlAuthorBooksRepository) Delete(ctx context.Context, id string) error {
	err := m.db.WithContext(ctx).Where("author_id = ?", id).Delete(&domain.AuthorBooks{}).Error
	return err
}
