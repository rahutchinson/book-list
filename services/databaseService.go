package services

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"

	"github.com/rahutchinson/book-list/models"
)

func ReadBookRows(conn *pgx.Conn) models.Books {
	rows, err := conn.Query(context.Background(), "SELECT isbn, name, author, type, description, cover, genre, tags, link FROM books")
	if err != nil {
		log.Fatal(err)
	}
	var bS models.Books
	defer rows.Close()
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ISBN, &book.Name, &book.Author, &book.Type, &book.Description, &book.Cover, &book.Genre, &book.Tags, &book.Link); err != nil {
			fmt.Println(err)
		}
		bS.Books = append(bS.Books, book)
	}
	return bS
}

func ReadRowsFeatured(conn *pgx.Conn) models.FeaturedBooks {
	rows, err := conn.Query(context.Background(), "SELECT isbn, current FROM featured")
	if err != nil {
		log.Fatal(err)
	}
	var bS models.FeaturedBooks
	defer rows.Close()
	for rows.Next() {
		var book models.FeaturedBook
		if err := rows.Scan(&book.ISBN, &book.Current); err != nil {
			fmt.Println(err)
		}
		bS.Featured = append(bS.Featured, book)
	}
	return bS
}