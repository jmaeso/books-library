package storage

import library "github.com/jmaeso/books-library/books-library"

type BooksStore interface {
	Insert(book library.Book) error
	DeleteForUser(pk int64, username string) error
	FindForUser(pk int64, username string) (*library.Book, error)
	GetAllSortedAndFilteredForUser(sortBy, filter, username string) (*[]library.Book, error)
}
