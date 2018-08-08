package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/jmaeso/books-library/books-library"
	"github.com/jmaeso/books-library/books-library/api"
	"github.com/jmaeso/books-library/books-library/storage/sqlite3"
	"github.com/jmaeso/books-library/pkg/classify"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/negroni"
	"gopkg.in/gorp.v2"
)

func initDB() (*sql.DB, *gorp.DbMap) {
	var err error

	db, err := sql.Open("sqlite3", "dev.db")
	if err != nil {
		log.Fatal(err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(library.Book{}, "books").SetKeys(true, "pk")
	dbmap.AddTableWithName(library.User{}, "users").SetKeys(false, "username")
	dbmap.CreateTablesIfNotExists()

	return db, dbmap
}

func main() {
	db, dbmap := initDB()

	classifyClient := classify.NewClient()

	sqlite3BooksStore := &sqlite3.BooksStore{
		DBMap: dbmap,
	}

	mux := mux.NewRouter()

	mux.HandleFunc("/", api.GetRootHandler(sqlite3BooksStore)).Methods("GET")

	mux.HandleFunc("/login", api.LoginHandler(dbmap))
	mux.HandleFunc("/logout", api.LogoutHandler())

	mux.HandleFunc("/books", api.GetFilteredBooksHandler(sqlite3BooksStore)).
		Methods("GET").
		Queries("filter", "{filter:all|fiction|nonfiction}")

	mux.HandleFunc("/books", api.GetSortedBooksHandler(sqlite3BooksStore)).
		Methods("GET").
		Queries("sortBy", "{sortBy:title|author|classification}")

	mux.HandleFunc("/books", api.CreateBooksHandler(sqlite3BooksStore, classifyClient)).Methods("PUT")

	mux.HandleFunc("/books/{pk}", api.DeleteBooksHandler(sqlite3BooksStore)).Methods("DELETE")

	mux.HandleFunc("/search", api.PostSearchHandler(classifyClient)).Methods("POST")

	n := negroni.Classic()
	n.Use(sessions.Sessions("books-library", cookiestore.New([]byte(os.Getenv("BOOKS_LIBRARY_COOKIES_KEY")))))
	n.Use(negroni.HandlerFunc(api.VerifyDB(db)))
	n.Use(negroni.HandlerFunc(api.VerifyUser(dbmap)))
	n.UseHandler(mux)

	n.Run(":8080")
}
