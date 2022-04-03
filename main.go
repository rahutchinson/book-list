package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/jackc/pgx/v4"

	models "github.com/rahutchinson/book-list/models"
	services "github.com/rahutchinson/book-list/services"
)

var (
	httpAddr     = flag.String("http", defaultAddr(), "http listen address")
	postKey      = os.Getenv("POST_KEY")
	dbConnection = os.Getenv("DB_CON_STRING")
	index        *template.Template
)

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
	services.ReadBookRows(dbConn)

	flag.Parse()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/books", bookHandler)
	http.HandleFunc("/featured", featuredHandler)
	fs := http.FileServer(http.Dir("./js/"))
	http.Handle("/js/", http.StripPrefix("/js", fs))

	log.Print("Running at address ", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	params := models.IndexParams{
		Host: req.Host,
	}
	w.Header().Set("Cache-Control", "no-cache")
	index.Execute(w, params)
}

func featuredHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		s := services.ReadRowsFeatured(dbConn)
		err := json.NewEncoder(w).Encode(s)
		if err != nil {
			return
		}
	case http.MethodPost:
		var b models.PostFeatured
		err := json.NewDecoder(req.Body).Decode(&b)
		if err != nil {
			http.Error(w, "Bad POST", 400)
		}
		if b.Key == postKey {
			repeat := services.AddToFeatured(b.FeaturedBook, dbConn)
			if !repeat {
				http.Error(w, "repeat or bad book", http.StatusConflict)
			}
		} else {
			http.Error(w, "F off", 404)
		}
	case http.MethodPut:
		var b models.PostFeatured
		err := json.NewDecoder(req.Body).Decode(&b)
		if err != nil {
			http.Error(w, "Bad POST", 400)
		}
		if b.Key == postKey {
			repeat := services.UpdateFeatured(b.FeaturedBook, dbConn)
			if !repeat {
				http.Error(w, "repeat or bad book", http.StatusConflict)
			}
		} else {
			http.Error(w, "F off", 404)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func bookHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		s := services.ReadBookRows(dbConn)
		err := json.NewEncoder(w).Encode(s)
		if err != nil {
			return
		}
	case http.MethodPost:
		var b models.PostBook
		err := json.NewDecoder(req.Body).Decode(&b)
		if err != nil {
			http.Error(w, "Bad POST", 400)
		}
		if b.Key == postKey {
			repeat := services.AddToBook(b.Book, dbConn)
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
	}

	rand.Seed(time.Now().UnixNano())
}
