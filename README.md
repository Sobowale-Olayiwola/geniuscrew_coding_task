# geniuscrew_coding_task
Create two services for managing books and authors, you can use docker and one database
or you can use a memory database (eg. with struct, map,slice).

### First service for books you must implement API like REST:
* Create a book with multiple authors, books have fields: id, title, description, authors,
publication date and other fields you can add fields for describing a book inside a
database.
* Get a single book with id
* Update a book inside db, manage association with authors
* Delete a book with id
* Search books for title or description or authorId and implement join with authors

### Second service for authors you must implement API like REST:
* Create an author, authors have fields:id, name, surname, books published and others
fields you can add fields for describing authors inside a database.
* Delete an author, when delete an author you must delete all books of author
* Get a single author with id
* Delete an author with id
* Update author with id
* Search authors and implement join result with books

## Endpoints available
baseurl = localhost 
port = 8080

### Creates an author
* POST 
    * /api/v1/authors
### Fetches an author with a specific id
* GET 
    * /api/v1/authors/:id
### Fetches an author with a filter based on a field e.g name,surname,email
* GET 
    * /api/v1/authors/filter

### Updates an author with a specific id
* PUT 
    * /api/v1/authors/:id

### Deletes
 an author with a specific id
* DELETE 
    * /api/v1/authors/:id

### Creates a book
* POST 
    * /api/v1/books
### Fetches a book with a specific id
* GET 
    * /api/v1/books/:id
### Fetches a book with a filter based on a field e.g title, description
* GET 
    * /api/v1/books/filter

### Updates a book with a specific id
* PUT 
    * /api/v1/books/:id

### Deletes a book with a specific id
* DELETE 
    * /api/v1/books/:id
## How to run generate executable
* go mod download
* cd cmd/api
* go build main.go data_sources.go injection.go
* ./main
## How to generate test report for each test file
### coverage 
* cd author/service
* go test -cover

The command generates a coverage of 100%

### Test
* cd author/service
* go test -v

### coverage
* cd book/service
* go test -cover

The command generates a coverage of 100%

### Test
* cd author/service
* go test -v