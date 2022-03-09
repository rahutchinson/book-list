package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/jackc/pgx/v4"
)

var (
	httpAddr     = flag.String("http", defaultAddr(), "http listen address")
	postKey      = os.Getenv("POST_KEY")
	dbConnection = os.Getenv("DB_CON_STRING")
	index        *template.Template
)

type books struct {
	Books []book `json:"books"`
}

type book struct {
	ISBN        string `json:"isbn"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	Genre       string `json:"genre"`
	Tags        string `json:"tags"`
	Link        string `json:"link"`
}

type indexParams struct {
	Host string
}

type postBook struct {
	Book book
	Key  string
}

var dbConn *pgx.Conn

func main() {
	config, err := pgx.ParseConfig(os.ExpandEnv(dbConnection))
	config.Database = "defaultdb"
	if err != nil {
		log.Fatal("error configuring the database: ", err)
	}
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	dbConn = conn
	readRows(dbConn)

	flag.Parse()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/books", bookHandler)
	fs := http.FileServer(http.Dir("./js/"))
	http.Handle("/js/", http.StripPrefix("/js", fs))

	log.Print("Running at address ", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

func readRows(conn *pgx.Conn) books {
	rows, err := conn.Query(context.Background(), "SELECT isbn, name, author, type, description, cover, genre, tags, link FROM books")
	if err != nil {
		log.Fatal(err)
	}
	var bS books
	defer rows.Close()
	for rows.Next() {
		var book book
		if err := rows.Scan(&book.ISBN, &book.Name, &book.Author, &book.Type, &book.Description, &book.Cover, &book.Genre, &book.Tags, &book.Link); err != nil {
			fmt.Println(err)
		}
		bS.Books = append(bS.Books, book)
	}
	return bS
}

func insertRows(ctx context.Context, tx pgx.Tx, bookToAdd book) error {
	// Insert four rows into the "accounts" table.
	log.Println("Creating new rows...")
	if _, err := tx.Exec(ctx,
		"INSERT INTO public.books (isbn, name, author, type, description, cover, genre, tags, link) VALUES ($1, $2, $3, $4,$5, $6, $7, $8, $9)", bookToAdd.ISBN, bookToAdd.Name, bookToAdd.Author, bookToAdd.Type, bookToAdd.Description, bookToAdd.Cover, bookToAdd.Genre, bookToAdd.Tags, bookToAdd.Link); err != nil {
		return err
	}
	return nil
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	params := indexParams{
		Host: req.Host,
	}
	w.Header().Set("Cache-Control", "no-cache")
	index.Execute(w, params)
}

func bookHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		s := readRows(dbConn)
		err := json.NewEncoder(w).Encode(s)
		if err != nil {
			return
		}
	case http.MethodPost:
		var b postBook
		err := json.NewDecoder(req.Body).Decode(&b)
		if err != nil {
			http.Error(w, "Bad POST", 400)
		}
		if b.Key == postKey {
			repeat := addToBook(b.Book)
			if !repeat {
				http.Error(w, "repeat or bad book", http.StatusConflict)
			}
		} else {
			http.Error(w, "F off", 404)
		}

	case http.MethodPut:
		// Update an existing record.
	case http.MethodDelete:
		// Remove the record.
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func addToBook(b book) bool {
	if b.ISBN == "" || b.Link == "" || b.Name == "" {
		return false
	}
	err := crdbpgx.ExecuteTx(context.Background(), dbConn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return insertRows(context.Background(), tx, b)
	})
	if err == nil {
		log.Println("New rows created.")
	} else {
		log.Fatal("error: ", err)
	}
	return true
}

func defaultAddr() string {
	port := os.Getenv("PORT")
	if port != "" {
		return ":" + port
	}

	return ":8080"
}

func init() {
	var err error

	// Parse optional on-disk index file.
	if index, err = template.ParseFiles("./index.html"); err != nil {
		log.Println(err)
		log.Println("Using default template")
		index = template.Must(template.New("index").Parse(indexHtml))
	}

	rand.Seed(time.Now().UnixNano())
}

var indexHtml = `
`
