package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	models "github.com/rahutchinson/book-list/models"
)

// SQLBook represents the structure of the source SQL table
type SQLBook struct {
	ISBN        string `db:"isbn"`
	Name        string `db:"name"`
	Author      string `db:"author"`
	Type        string `db:"type"`
	Description string `db:"description"`
	Cover       string `db:"cover"`
	Genre       string `db:"genre"`
	Tags        string `db:"tags"`
	Link        string `db:"link"`
}

func main() {
	// Database connection parameters
	dbUser := getEnv("DB_USER", "script")
	dbPassword := getEnv("DB_PASSWORD", "vBbmP0zJJQsz5Q_0Qzdymw")
	dbHost := getEnv("DB_HOST", "loyal-efreet-5669.5xj.gcp-us-central1.cockroachlabs.cloud")
	dbPort := getEnv("DB_PORT", "26257")
	dbName := getEnv("DB_NAME", "defaultdb")
	outputFile := getEnv("OUTPUT_FILE", "books.json")

	// Connect to CockroachDB
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	
	// Alternative: Use the exact connection string format
	// dsn := "postgresql://ryan:vBbmP0zJJQsz5Q_0Qzdymw@loyal-efreet-5669.5xj.gcp-us-central1.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database successfully")

	// Query books from SQL table
	rows, err := db.Query("SELECT isbn, name, author, type, description, cover, genre, tags, link FROM books")
	if err != nil {
		log.Fatalf("Failed to query books: %v", err)
	}
	defer rows.Close()

	var books models.Books
	var count int

	for rows.Next() {
		var sqlBook SQLBook
		err := rows.Scan(&sqlBook.ISBN, &sqlBook.Name, &sqlBook.Author, &sqlBook.Type,
			&sqlBook.Description, &sqlBook.Cover, &sqlBook.Genre, &sqlBook.Tags, &sqlBook.Link)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		// Convert SQL book to application book model
		book := convertSQLBookToModel(sqlBook)
		books.Books = append(books.Books, book)
		count++
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating over rows: %v", err)
	}

	// Write to JSON file
	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile(outputFile, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Printf("Successfully migrated %d books to %s\n", count, outputFile)
	fmt.Printf("Migration completed at %s\n", time.Now().Format(time.RFC3339))
}

func convertSQLBookToModel(sqlBook SQLBook) models.Book {
	// Generate unique ID
	id := fmt.Sprintf("%d_%s", time.Now().UnixNano(), sqlBook.ISBN)

	// Map type to application enum
	bookType := models.Physical // default
	switch sqlBook.Type {
	case "physical":
		bookType = models.Physical
	case "kindle":
		bookType = models.Kindle
	case "audible":
		bookType = models.Audible
	case "ebook":
		bookType = models.Ebook
	}

	// Convert tags string to array
	var tags []string
	if sqlBook.Tags != "" {
		tags = []string{sqlBook.Tags}
	}

	return models.Book{
		ID:          id,
		ISBN:        sqlBook.ISBN,
		Name:        sqlBook.Name,
		Author:      sqlBook.Author,
		Type:        bookType,
		Description: sqlBook.Description,
		Cover:       sqlBook.Cover,
		Genre:       sqlBook.Genre,
		Tags:        tags,
		Link:        sqlBook.Link,
		Status:      models.Unread, // Default status
		Rating:      0,             // Default rating
		Pages:       0,             // Default pages
		Duration:    "",            // Default duration
		Publisher:   "",            // Default publisher
		Published:   time.Time{},   // Default published date
		Added:       time.Now(),    // Current timestamp
		Started:     time.Time{},   // Default started date
		Finished:    time.Time{},   // Default finished date
		Notes:       "",            // Default notes
		Series:      "",            // Default series
		SeriesOrder: 0,             // Default series order
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
