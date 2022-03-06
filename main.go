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

	index *template.Template
)

type books struct {
	Books []book `json:"books"`
}

type book struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	Link        string `json:"link"`
}

type indexParams struct {
	Books books
	Host  string
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
		Books: readBooks(),
		Host:  req.Host,
	}
	w.Header().Set("Cache-Control", "no-cache")
	index.Execute(w, params)
}

func bookHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s := readBooks()
	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		return
	}
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
