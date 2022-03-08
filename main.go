package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"
)

var (
	httpAddr = flag.String("http", defaultAddr(), "http listen address")
	postKey  = os.Getenv("POST_KEY")
	index    *template.Template
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

func main() {
	flag.Parse()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/books", bookHandler)
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
		s := readBooks()
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
	exsistingBooks := readBooks().Books

	var repeat bool
	for _, book := range exsistingBooks {
		repeat = b.ISBN == book.ISBN
	}
	if repeat {
		return false
	}
	booksToWrite := append(exsistingBooks, b)
	writeBooks(books{Books: booksToWrite})
	fmt.Println("write to file")
	return true
}

func writeBooks(toWrite books) {
	file, _ := json.MarshalIndent(toWrite, "", " ")
	_ = ioutil.WriteFile("books.json", file, 0644)
}

func readBooks() books {
	// Open our jsonFile
	jsonFile, err := os.Open("./books.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened books.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var allbooks books

	err = json.Unmarshal(byteValue, &allbooks)
	if err != nil {
		fmt.Println("error")
		return books{}
	}
	return allbooks
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
