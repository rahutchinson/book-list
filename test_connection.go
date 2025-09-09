package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Test connection string - try different SSL modes
	
	fmt.Println("Attempting to connect to CockroachDB...")
	fmt.Printf("Connection string: %s\n", dsn)
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to database!")

	// Test if the books table exists
	var tableExists bool
	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'books')").Scan(&tableExists)
	if err != nil {
		log.Fatalf("Failed to check if books table exists: %v", err)
	}

	if tableExists {
		fmt.Println("Books table exists!")
		
		// Count books in the table
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
		if err != nil {
			log.Fatalf("Failed to count books: %v", err)
		}
		fmt.Printf("Found %d books in the table\n", count)
		
		// Show a sample book
		var isbn, name, author string
		err = db.QueryRow("SELECT isbn, name, author FROM books LIMIT 1").Scan(&isbn, &name, &author)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("No books found in the table")
			} else {
				log.Fatalf("Failed to get sample book: %v", err)
			}
		} else {
			fmt.Printf("Sample book: ISBN=%s, Name=%s, Author=%s\n", isbn, name, author)
		}
	} else {
		fmt.Println("Books table does not exist!")
	}
}
