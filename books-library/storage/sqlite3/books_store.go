package sqlite3

import (
	library "github.com/jmaeso/books-library/books-library"
	"gopkg.in/gorp.v2"
)

const (
	selectBookByPKAndUsernameStmt = "SELECT pk, title, author, classification, id, username FROM books WHERE pk=? AND username=?"
	deleteBookByPKAndUsernameStmt = "DELETE FROM books WHERE pk=? AND username=?"
	selectAllByUsernameStmt       = "SELECT pk, title, author, classification, id, username FROM books WHERE username=?"
)

type BooksStore struct {
	DBMap *gorp.DbMap
}

func (bs *BooksStore) Insert(book library.Book) error {
	return bs.DBMap.Insert(&book)
}

func (bs *BooksStore) DeleteForUser(pk int64, username string) error {
	if _, err := bs.DBMap.Exec(deleteBookByPKAndUsernameStmt, pk, username); err != nil {
		return err
	}

	return nil
}

func (bs *BooksStore) FindForUser(pk int64, username string) (*library.Book, error) {
	var b library.Book
	if err := bs.DBMap.SelectOne(&b, selectBookByPKAndUsernameStmt, pk, username); err != nil {
		return nil, err
	}

	return &b, nil
}

func (bs *BooksStore) GetAllSortedAndFilteredForUser(sortBy, filter, username string) (*[]library.Book, error) {
	var books []library.Book

	var filterStmt string
	if filter == "fiction" {
		filterStmt = " and classification between '800' and '900'"
	} else if filter == "nonfiction" {
		filterStmt = " and classification not between '800' and '900'"
	}

	var orderByStmt string
	if sortBy == "" {
		orderByStmt = " ORDER BY pk"
	} else {
		orderByStmt = " ORDER BY " + sortBy
	}

	if _, err := bs.DBMap.Select(&books, selectAllByUsernameStmt+filterStmt+orderByStmt, username); err != nil {
		return nil, err
	}

	return &books, nil
}
