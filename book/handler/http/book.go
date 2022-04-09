package http

import (
	"context"
	"errors"
	"geniuscrew/domain"
	"geniuscrew/internal/appvalidator"
	"geniuscrew/internal/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	BookService domain.BookService
}

func NewBookHandler(router *gin.Engine, p domain.BookService) {
	handler := &BookHandler{
		BookService: p,
	}
	api := router.Group("/api/v1")
	api.POST("/books", handler.CreateBook)
	api.GET("/books/:id", handler.GetBookByID)
	api.GET("/books/filter", handler.GetByFilter)
	api.PUT("/books/:id", handler.UpdateBookByID)
	api.DELETE("/books/:id", handler.DeleteBookByID)
}

func (p *BookHandler) CreateBook(c *gin.Context) {
	var input struct {
		Title             string `json:"title" validate:"gte=0,lte=500,required"`
		Description       string `json:"description" validate:"gte=0,lte=500,required"`
		ISBN              string `json:"ISBN" validate:"gte=0,lte=14,required"`
		PublishingCompany string `json:"publishing_company" validate:"gte=0,lte=50,required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inputErr := appvalidator.InputValidator(input)
	if inputErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": inputErr})
		return
	}

	var ctx = context.TODO()
	var book domain.Book
	book.Title = input.Title
	book.Description = input.Description
	book.ISBN = input.ISBN
	book.PublishingCompany = input.PublishingCompany

	err := p.BookService.Create(ctx, &book)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicateRecord):
			c.JSON(http.StatusConflict, gin.H{"error": "book exists"})
			return
		case errors.Is(err, domain.ErrBookNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"payload": book})
}

func (p *BookHandler) GetBookByID(c *gin.Context) {
	id := c.Param("id")
	err := appvalidator.IsIDValid(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	var ctx = context.TODO()
	book, err := p.BookService.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusFound, gin.H{"payload": book})
}

func (p *BookHandler) UpdateBookByID(c *gin.Context) {
	id := c.Param("id")
	err := appvalidator.IsIDValid(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	var input struct {
		Title             string `json:"title" validate:"isdefault|gte=0,lte=500"`
		Description       string `json:"description" validate:"isdefault|gte=0,lte=500"`
		ISBN              string `json:"ISBN" validate:"isdefault|gte=0,lte=14"`
		PublishingCompany string `json:"publishing_company" validate:"isdefault|gte=0, lte=50"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inputErr := appvalidator.InputValidator(input)
	if inputErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": inputErr})
		return
	}
	var ctx = context.TODO()
	book, err := p.BookService.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	var updatedBook domain.Book
	updatedBook.Title = input.Title
	updatedBook.Description = input.Description
	updatedBook.ISBN = input.ISBN
	updatedBook.PublishingCompany = input.PublishingCompany
	err = p.BookService.Update(ctx, id, &book, updatedBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (p *BookHandler) DeleteBookByID(c *gin.Context) {
	id := c.Param("id")
	err := appvalidator.IsIDValid(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	var ctx = context.TODO()
	var book domain.Book
	err = p.BookService.Delete(ctx, id, &book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func (p *BookHandler) GetByFilter(c *gin.Context) {
	filter := c.Query("field")
	filterValue := c.Query("value")
	filterSafeList := []string{"title", "description", "authorId"}
	if !helpers.In(filter, filterSafeList...) {
		message := make(map[string][]string)
		message["filter_fields_allowed"] = filterSafeList
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": message})
		return
	}
	var ctx = context.TODO()

	book, err := p.BookService.GetByFilter(ctx, filter, filterValue)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBookNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusFound, gin.H{"payload": book})
}
