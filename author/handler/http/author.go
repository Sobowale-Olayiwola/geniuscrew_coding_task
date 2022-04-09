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

type AuthorHandler struct {
	AuthorService domain.AuthorService
}

func NewAuthorHandler(router *gin.Engine, as domain.AuthorService) {
	handler := &AuthorHandler{
		AuthorService: as,
	}
	api := router.Group("/api/v1")
	api.POST("/authors", handler.CreateAuthor)
	api.GET("/authors/:id", handler.GetAuthorByID)
	api.GET("/authors/filter", handler.GetByFilter)
	api.PUT("/authors/:id", handler.UpdateAuthorByID)
	api.DELETE("/authors/:id", handler.DeleteAuthorByID)
}

func (p *AuthorHandler) CreateAuthor(c *gin.Context) {
	var input struct {
		Name           string   `json:"name" validate:"gte=0,lte=500,required"`
		Surname        string   `json:"surname" validate:"gte=0,lte=500,required"`
		Email          string   `json:"email" validate:"email,required"`
		BooksPublished []string `json:"books_published" validate:"min=1,dive,required"`
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
	var author domain.Author
	author.Name = input.Name
	author.Surname = input.Surname
	author.Email = input.Email
	err := p.AuthorService.Create(ctx, input.BooksPublished, &author)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicateRecord):
			c.JSON(http.StatusConflict, gin.H{"error": "Author Email is registered"})
			return
		case errors.Is(err, domain.ErrBookNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"payload": author})
}

func (p *AuthorHandler) GetAuthorByID(c *gin.Context) {
	id := c.Param("id")
	err := appvalidator.IsIDValid(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	var ctx = context.TODO()
	author, err := p.AuthorService.Get(ctx, id)
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
	c.JSON(http.StatusFound, gin.H{"payload": author})
}

func (p *AuthorHandler) GetByFilter(c *gin.Context) {
	filter := c.Query("field")
	filterValue := c.Query("value")
	filterSafeList := []string{"name", "surname", "email"}
	if !helpers.In(filter, filterSafeList...) {
		message := make(map[string][]string)
		message["filter_fields_allowed"] = filterSafeList
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": message})
		return
	}
	var ctx = context.TODO()

	author, err := p.AuthorService.GetByFilter(ctx, filter, filterValue)
	if len(author) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no author matches"})
		return
	}
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
	c.JSON(http.StatusFound, gin.H{"payload": author})
}

func (p *AuthorHandler) UpdateAuthorByID(c *gin.Context) {
	id := c.Param("id")
	err := appvalidator.IsIDValid(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	var input struct {
		Name           string   `json:"name" validate:"isdefault|gte=0,lte=500"`
		Surname        string   `json:"surname" validate:"isdefault|gte=0,lte=500"`
		Email          string   `json:"email" validate:"isdefault|email"`
		BooksPublished []string `json:"books_published" validate:"isdefault|min=1,dive,required"`
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
	author, err := p.AuthorService.Get(ctx, id)
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
	var updatedAuthor domain.Author
	updatedAuthor.Email = input.Email
	updatedAuthor.Name = input.Name
	updatedAuthor.Surname = input.Surname
	err = p.AuthorService.Update(ctx, id, &author, updatedAuthor, input.BooksPublished)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "author profile updated",
	})
}

func (p *AuthorHandler) DeleteAuthorByID(c *gin.Context) {
	id := c.Param("id")
	err := appvalidator.IsIDValid(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	var ctx = context.TODO()
	var author domain.Author
	err = p.AuthorService.Delete(ctx, id, &author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Author deleted successfully"})
}
