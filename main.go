package main

import (
	"database/sql"
	"log"

	"github.com/jmaeso/books-library/pkg/classify"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	gmux "github.com/gorilla/mux"
	library "github.com/jmaeso/books-library/books-library"
	"github.com/jmaeso/books-library/books-library/api"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/negroni"
	"gopkg.in/gorp.v2"
)

var db *sql.DB
var dbmap *gorp.DbMap

func initDB() {
	var err error

	db, err = sql.Open("sqlite3", "dev.db")
	if err != nil {
		log.Fatal(err)
	}

	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(library.Book{}, "books").SetKeys(true, "pk")
	dbmap.AddTableWithName(library.User{}, "users").SetKeys(false, "username")
	dbmap.CreateTablesIfNotExists()
}

func main() {
	initDB()

	classifyClient := classify.NewClient()

	mux := gmux.NewRouter()

	mux.HandleFunc("/", api.GetRootHandler(dbmap)).Methods("GET")

	mux.HandleFunc("/login", api.LoginHandler(dbmap))
	mux.HandleFunc("/logout", api.LogoutHandler())

	mux.HandleFunc("/books", api.GetFilteredBooksHandler(dbmap)).
		Methods("GET").
		Queries("filter", "{filter:all|fiction|nonfiction}")

	mux.HandleFunc("/books", api.GetSortedBooksHandler(dbmap)).
		Methods("GET").
		Queries("sortBy", "{sortBy:title|author|classification}")

	mux.HandleFunc("/books", api.CreateBooksHandler(dbmap, classifyClient)).Methods("PUT")

	mux.HandleFunc("/books/{pk}", api.DeleteBooksHandler(dbmap)).Methods("DELETE")

	mux.HandleFunc("/search", api.PostSearchHandler(classifyClient)).Methods("POST")

	n := negroni.Classic()
	n.Use(sessions.Sessions("go-for-web-dev", cookiestore.New([]byte("my-secret-123"))))
	n.Use(negroni.HandlerFunc(api.VerifyDB(db)))
	n.Use(negroni.HandlerFunc(api.VerifyUser(dbmap)))
	n.UseHandler(mux)

	n.Run(":8080")
}
