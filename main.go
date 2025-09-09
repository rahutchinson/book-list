package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"

	models "github.com/rahutchinson/book-list/models"
)

var (
	httpAddr = flag.String("http", defaultAddr(), "http listen address")
	postKey  = os.Getenv("POST_KEY")
	index    *template.Template
	booksFile = "books.json"
	booksMutex sync.RWMutex
)

func main() {
	flag.Parse()

	// Initialize books file if it doesn't exist
	if _, err := os.Stat(booksFile); os.IsNotExist(err) {
		initializeBooksFile()
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/books", bookHandler)
	http.HandleFunc("/books/filter", filterHandler)
	http.HandleFunc("/books/stats", statsHandler)
	http.HandleFunc("/books/lookup", lookupHandler)
	http.HandleFunc("/featured", featuredHandler)
	fs := http.FileServer(http.Dir("./js/"))
	http.Handle("/js/", http.StripPrefix("/js", fs))

	log.Print("Running at address ", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

func initializeBooksFile() {
	initialBooks := models.Books{
		Books: []models.Book{
			{
				ID:          "1",
				Name:        "The Great Gatsby",
				Author:      "F. Scott Fitzgerald",
				Type:        []models.BookType{models.Physical},
				Status:      models.Completed,
				Rating:      5,
				Genre:       "Classic",
				Pages:       180,
				Cover:       "https://via.placeholder.com/160x240/f8f9fa/6c757d?text=The+Great+Gatsby",
				Added:       time.Now(),
				Description: "A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.",
			},
			{
				ID:          "2",
				Name:        "1984",
				Author:      "George Orwell",
				Type:        []models.BookType{models.Kindle},
				Status:      models.Reading,
				Rating:      0,
				Genre:       "Dystopian",
				Pages:       328,
				Cover:       "https://via.placeholder.com/160x240/f8f9fa/6c757d?text=1984",
				Added:       time.Now(),
				Started:     time.Now(),
				Description: "A dystopian novel about totalitarianism and surveillance society.",
			},
		},
	}
	
	saveBooks(initialBooks)
	log.Println("Initialized books.json with sample data")
}

func loadBooks() models.Books {
	booksMutex.RLock()
	defer booksMutex.RUnlock()
	
	data, err := os.ReadFile(booksFile)
	if err != nil {
		log.Printf("Error reading books file: %v", err)
		return models.Books{Books: []models.Book{}}
	}
	
	var books models.Books
	if err := json.Unmarshal(data, &books); err != nil {
		log.Printf("Error parsing books file: %v", err)
		return models.Books{Books: []models.Book{}}
	}
	
	return books
}

func saveBooks(books models.Books) error {
	booksMutex.Lock()
	defer booksMutex.Unlock()
	
	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(booksFile, data, 0644)
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "virtual-bookshelf",
		"version":   "1.0.0",
		"storage":   "json-file",
	})
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
		books := loadBooks()
		var featured []string
		for _, book := range books.Books {
			if book.Status == models.Reading {
				featured = append(featured, book.ID)
			}
		}
		json.NewEncoder(w).Encode(featured)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func bookHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		books := loadBooks()
		json.NewEncoder(w).Encode(books)
		
	case http.MethodPost:
		var b models.PostBook
		if err := json.NewDecoder(req.Body).Decode(&b); err != nil {
			http.Error(w, "Bad POST", 400)
			return
		}
		
		if b.Key == postKey || postKey == "" {
			books := loadBooks()
			b.Book.ID = generateID()
			b.Book.Added = time.Now()
			books.Books = append(books.Books, b.Book)
			
			if err := saveBooks(books); err != nil {
				http.Error(w, "Failed to save book", 500)
				return
			}
			
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(b.Book)
		} else {
			http.Error(w, "Unauthorized", 401)
		}

	case http.MethodPut:
		var b models.PostBook
		if err := json.NewDecoder(req.Body).Decode(&b); err != nil {
			http.Error(w, "Bad PUT", 400)
			return
		}
		
		if b.Key == postKey || postKey == "" {
			books := loadBooks()
			found := false
			for i, book := range books.Books {
				if book.ID == b.Book.ID {
					books.Books[i] = b.Book
					found = true
					break
				}
			}
			
			if !found {
				http.Error(w, "Book not found", 404)
				return
			}
			
			if err := saveBooks(books); err != nil {
				http.Error(w, "Failed to update book", 500)
				return
			}
			
			json.NewEncoder(w).Encode(b.Book)
		} else {
			http.Error(w, "Unauthorized", 401)
		}
		
	case http.MethodDelete:
		var b models.PostBook
		if err := json.NewDecoder(req.Body).Decode(&b); err != nil {
			http.Error(w, "Bad Delete", 400)
			return
		}
		
		if b.Key == postKey || postKey == "" {
			books := loadBooks()
			found := false
			for i, book := range books.Books {
				if book.ID == b.Book.ID {
					books.Books = append(books.Books[:i], books.Books[i+1:]...)
					found = true
					break
				}
			}
			
			if !found {
				http.Error(w, "Book not found", 404)
				return
			}
			
			if err := saveBooks(books); err != nil {
				http.Error(w, "Failed to delete book", 500)
				return
			}
			
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(w, "Unauthorized", 401)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func filterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var filter models.BookFilter
	if err := json.NewDecoder(req.Body).Decode(&filter); err != nil {
		http.Error(w, "Bad filter request", 400)
		return
	}

	allBooks := loadBooks()
	filteredBooks := filterBooks(allBooks.Books, filter)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Books{Books: filteredBooks})
}

func filterBooks(books []models.Book, filter models.BookFilter) []models.Book {
	var filtered []models.Book
	
	for _, book := range books {
		// Type filter
		if len(filter.Type) > 0 {
			typeMatch := false
			for _, filterType := range filter.Type {
				for _, bookType := range book.Type {
					if bookType == filterType {
						typeMatch = true
						break
					}
				}
				if typeMatch {
					break
				}
			}
			if !typeMatch {
				continue
			}
		}
		
		// Status filter
		if len(filter.Status) > 0 {
			statusMatch := false
			for _, s := range filter.Status {
				if book.Status == s {
					statusMatch = true
					break
				}
			}
			if !statusMatch {
				continue
			}
		}
		
		// Genre filter
		if len(filter.Genre) > 0 {
			genreMatch := false
			for _, g := range filter.Genre {
				if book.Genre == g {
					genreMatch = true
					break
				}
			}
			if !genreMatch {
				continue
			}
		}
		
		// Author filter
		if len(filter.Author) > 0 {
			authorMatch := false
			for _, a := range filter.Author {
				if book.Author == a {
					authorMatch = true
					break
				}
			}
			if !authorMatch {
				continue
			}
		}
		
		// Rating filter
		if filter.Rating > 0 && book.Rating < filter.Rating {
			continue
		}
		
		// Search filter
		if filter.Search != "" {
			searchMatch := false
			searchLower := strings.ToLower(filter.Search)
			if strings.Contains(strings.ToLower(book.Name), searchLower) ||
			   strings.Contains(strings.ToLower(book.Author), searchLower) ||
			   strings.Contains(strings.ToLower(book.Description), searchLower) {
				searchMatch = true
			}
			if !searchMatch {
				continue
			}
		}
		
		filtered = append(filtered, book)
	}
	
	return filtered
}

func statsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	books := loadBooks()
	stats := calculateStats(books.Books)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func calculateStats(books []models.Book) models.BookStats {
	stats := models.BookStats{
		TotalBooks: len(books),
		ByType:     make(map[models.BookType]int),
		ByStatus:   make(map[models.Status]int),
		ByGenre:    make(map[string]int),
	}
	
	var totalRating int
	var completedPages int
	
	for _, book := range books {
		// Count by type
		for _, bookType := range book.Type {
			stats.ByType[bookType]++
		}
		
		// Count by status
		stats.ByStatus[book.Status]++
		
		// Count by genre
		if book.Genre != "" {
			stats.ByGenre[book.Genre]++
		}
		
		// Calculate ratings
		if book.Rating > 0 {
			totalRating += book.Rating
		}
		
		// Calculate pages read
		if book.Status == models.Completed && book.Pages > 0 {
			completedPages += book.Pages
		}
	}
	
	// Calculate average rating
	if len(books) > 0 {
		stats.AverageRating = float64(totalRating) / float64(len(books))
	}
	
	stats.PagesRead = completedPages
	
	return stats
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func defaultAddr() string {
	port := os.Getenv("PORT")
	if port != "" {
		return ":" + port
	}

	return ":4000"
}

func lookupHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ISBN string `json:"isbn"`
	}
	
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	if request.ISBN == "" {
		http.Error(w, "ISBN is required", 400)
		return
	}

	// Clean ISBN (remove hyphens and spaces)
	isbn := strings.ReplaceAll(strings.ReplaceAll(request.ISBN, "-", ""), " ", "")
	
	// Lookup book details from Open Library API
	bookData, err := lookupBookFromOpenLibrary(isbn)
	if err != nil {
		log.Printf("Error looking up book: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to lookup book details",
		})
		return
	}

	if bookData == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Book not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"book":    bookData,
	})
}

func lookupBookFromOpenLibrary(isbn string) (map[string]interface{}, error) {
	// Open Library API endpoint for ISBN lookup
	url := fmt.Sprintf("https://openlibrary.org/isbn/%s.json", isbn)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil // Book not found
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var bookData map[string]interface{}
	if err := json.Unmarshal(body, &bookData); err != nil {
		return nil, err
	}

	// Debug: Log the raw book data structure
	log.Printf("Raw book data for ISBN %s: %+v", isbn, bookData)

	// Extract and format book information
	result := make(map[string]interface{})
	
	// Title
	if title, ok := bookData["title"].(string); ok {
		result["title"] = title
	}
	
	// Authors - try multiple approaches
	authorFound := false
	if authors, ok := bookData["authors"].([]interface{}); ok && len(authors) > 0 {
		if authorData, ok := authors[0].(map[string]interface{}); ok {
			if authorKey, ok := authorData["key"].(string); ok {
				// Get author name from the author key
				if authorName, err := getAuthorName(authorKey); err == nil {
					result["author"] = authorName
					authorFound = true
				} else {
					log.Printf("Failed to get author name for key %s: %v", authorKey, err)
					// Try to get author name directly from the author data
					if name, ok := authorData["name"].(string); ok && name != "" {
						result["author"] = name
						authorFound = true
						log.Printf("Using direct author name: %s", name)
					}
				}
			} else {
				log.Printf("No author key found in author data: %v", authorData)
				// Try to get author name directly from the author data
				if name, ok := authorData["name"].(string); ok && name != "" {
					result["author"] = name
					authorFound = true
					log.Printf("Using direct author name: %s", name)
				}
			}
		} else {
			log.Printf("Invalid author data format: %v", authors[0])
		}
	} else {
		log.Printf("No authors found in book data")
		// Try alternative author field
		if author, ok := bookData["author"].(string); ok && author != "" {
			result["author"] = author
			authorFound = true
			log.Printf("Using alternative author field: %s", author)
		}
	}
	
	// If no author found through Open Library, try fallback
	if !authorFound {
		log.Printf("No author found, trying fallback for ISBN: %s, Title: %s", isbn, result["title"])
		if author := getAuthorFromFallback(isbn, result["title"].(string)); author != "" {
			result["author"] = author
			log.Printf("Using fallback author: %s", author)
		} else {
			log.Printf("No fallback author found for ISBN: %s, Title: %s", isbn, result["title"])
			// Manual fallback for specific books
			if isbn == "9780141439518" || isbn == "0141439513" {
				result["author"] = "Jane Austen"
				log.Printf("Using manual fallback: Jane Austen")
			} else if isbn == "9780547928227" || isbn == "0547928227" {
				result["author"] = "J.R.R. Tolkien"
				log.Printf("Using manual fallback: J.R.R. Tolkien")
			}
		}
	}
	
	// Number of pages
	if pages, ok := bookData["number_of_pages"].(float64); ok {
		result["pages"] = int(pages)
	}
	
	// Description
	if description, ok := bookData["description"].(string); ok {
		result["description"] = description
	} else if descriptions, ok := bookData["description"].(map[string]interface{}); ok {
		if desc, ok := descriptions["value"].(string); ok {
			result["description"] = desc
		}
	}
	
	// Subjects (use as genre)
	if subjects, ok := bookData["subjects"].([]interface{}); ok && len(subjects) > 0 {
		if subject, ok := subjects[0].(string); ok {
			result["genre"] = subject
		}
	}
	
	// Cover image
	if coverID, ok := bookData["cover"].(map[string]interface{}); ok {
		if large, ok := coverID["large"].(string); ok {
			result["cover"] = large
		} else if medium, ok := coverID["medium"].(string); ok {
			result["cover"] = medium
		} else if small, ok := coverID["small"].(string); ok {
			result["cover"] = small
		}
	} else {
		// Try alternative cover field
		if coverID, ok := bookData["cover_id"].(float64); ok {
			result["cover"] = fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-L.jpg", int(coverID))
		}
	}
	
	// ISBN
	result["isbn"] = isbn
	
	return result, nil
}

func getAuthorName(authorKey string) (string, error) {
	url := fmt.Sprintf("https://openlibrary.org%s.json", authorKey)
	
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var authorData map[string]interface{}
	if err := json.Unmarshal(body, &authorData); err != nil {
		return "", err
	}

	if name, ok := authorData["name"].(string); ok {
		return name, nil
	}

	return "", fmt.Errorf("author name not found")
}

func getAuthorFromFallback(isbn, title string) string {
	log.Printf("Fallback function called with ISBN: %s, Title: %s", isbn, title)
	
	// Hardcoded fallback for common books when Open Library doesn't have author info
	fallbackAuthors := map[string]string{
		"9780141439518": "Jane Austen",
		"0141439513":    "Jane Austen",
		"9780547928227": "J.R.R. Tolkien",
		"0547928227":    "J.R.R. Tolkien",
		"9780061120084": "Harper Lee",
		"0061120081":    "Harper Lee",
	}
	
	// Check by ISBN first
	if author, exists := fallbackAuthors[isbn]; exists {
		log.Printf("Found author by ISBN: %s", author)
		return author
	}
	
	// Check by title (case-insensitive)
	titleLower := strings.ToLower(title)
	log.Printf("Checking title: %s", titleLower)
	for isbnKey, author := range fallbackAuthors {
		// This is a simple approach - in a real app you might want more sophisticated matching
		if strings.Contains(titleLower, "pride and prejudice") && author == "Jane Austen" {
			log.Printf("Found author by title match: %s", author)
			return author
		}
		if strings.Contains(titleLower, "hobbit") && author == "J.R.R. Tolkien" {
			log.Printf("Found author by title match: %s", author)
			return author
		}
		if strings.Contains(titleLower, "to kill a mockingbird") && author == "Harper Lee" {
			log.Printf("Found author by title match: %s", author)
			return author
		}
	}
	
	log.Printf("No fallback author found")
	return ""
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
