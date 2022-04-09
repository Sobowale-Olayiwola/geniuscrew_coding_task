package main

import (
	_mysqlAuthorRepo "geniuscrew/author/repository/mysql"
	_mysqlBookRepo "geniuscrew/book/repository/mysql"

	_authorService "geniuscrew/author/service"
	_bookService "geniuscrew/book/service"

	_authorHandler "geniuscrew/author/handler/http"
	_bookHandler "geniuscrew/book/handler/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func inject(d *DataSources) *gin.Engine {
	/*
	 * repository layer
	 */
	mysqlBookRepo := _mysqlBookRepo.NewMySqlBookRepository(d.MySQLDB)
	mysqlAuthorRepo := _mysqlAuthorRepo.NewMySqlAuthorRepository(d.MySQLDB)
	mysqlAuthorBooksRepo := _mysqlAuthorRepo.NewMySqlAuthorBooksRepository(d.MySQLDB)

	/*
	 * service layer
	 */
	bookService := _bookService.NewBookService(mysqlBookRepo)
	authorService := _authorService.NewAuthorService(mysqlAuthorRepo, mysqlAuthorBooksRepo, mysqlBookRepo)

	router := gin.Default()

	router.Use(cors.Default())
	/*
	 * handler layer
	 */
	_bookHandler.NewBookHandler(router, bookService)
	_authorHandler.NewAuthorHandler(router, authorService)

	return router
}
