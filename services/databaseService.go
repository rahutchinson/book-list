package services

import (
	"context"
	"fmt"
	"log"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
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

func AddToBook(b models.Book, dbConn *pgx.Conn) bool {
	if b.ISBN == "" || b.Link == "" || b.Name == "" {
		return false
	}
	err := crdbpgx.ExecuteTx(context.Background(), dbConn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return insertRowsBooks(context.Background(), tx, b)
	})
	if err == nil {
		log.Println("New rows created.")
	} else {
		log.Fatal("error: ", err)
	}
	return true
}

func AddToFeatured(b models.FeaturedBook, dbConn *pgx.Conn) bool {
	if b.ISBN == "" {
		return false
	}
	err := crdbpgx.ExecuteTx(context.Background(), dbConn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return insertRowsFeatrued(context.Background(), tx, b)
	})
	if err == nil {
		log.Println("New rows created.")
	} else {
		log.Println("error: ", err)
		return false
	}
	return true
}

func UpdateFeatured(b models.FeaturedBook, dbConn *pgx.Conn) bool {
	if b.ISBN == "" {
		return false
	}
	err := crdbpgx.ExecuteTx(context.Background(), dbConn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return updateRowFeatrued(context.Background(), tx, b)
	})
	if err == nil {
		log.Println("New rows created.")
	} else {
		log.Println("error: ", err)
		return false
	}
	return true
}

func insertRowsBooks(ctx context.Context, tx pgx.Tx, bookToAdd models.Book) error {
	// Insert four rows into the "accounts" table.
	log.Println("Creating new rows...")
	if _, err := tx.Exec(ctx,
		"INSERT INTO public.books (isbn, name, author, type, description, cover, genre, tags, link) VALUES ($1, $2, $3, $4,$5, $6, $7, $8, $9)", bookToAdd.ISBN, bookToAdd.Name, bookToAdd.Author, bookToAdd.Type, bookToAdd.Description, bookToAdd.Cover, bookToAdd.Genre, bookToAdd.Tags, bookToAdd.Link); err != nil {
		return err
	}
	return nil
}

func insertRowsFeatrued(ctx context.Context, tx pgx.Tx, bookToAdd models.FeaturedBook) error {
	// Insert four rows into the "accounts" table.
	log.Println("Creating new rows...")
	if _, err := tx.Exec(ctx,
		"INSERT INTO public.featured (isbn, current) VALUES ($1, $2)", bookToAdd.ISBN, bookToAdd.Current); err != nil {
		return err
	}
	return nil
}

func updateRowFeatrued(ctx context.Context, tx pgx.Tx, bookToAdd models.FeaturedBook) error {
	// Insert four rows into the "accounts" table.
	log.Println("Updating featured...")
	if _, err := tx.Exec(ctx,
		"UPDATE public.featured SET current = $2 WHERE isbn = $1", bookToAdd.ISBN, bookToAdd.Current); err != nil {
		return err
	}
	return nil
}

func ReadRowsFeatured(conn *pgx.Conn) []string {
	rows, err := conn.Query(context.Background(), "SELECT isbn, current FROM featured WHERE current")
	if err != nil {
		log.Fatal(err)
	}
	var bS []string
	defer rows.Close()
	for rows.Next() {
		var book models.FeaturedBook
		if err := rows.Scan(&book.ISBN, &book.Current); err != nil {
			fmt.Println(err)
		}
		bS = append(bS, book.ISBN)
	}
	return bS
}
