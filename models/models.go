package models

import "time"

type Books struct {
	Books []Book `json:"books"`
}

type FeaturedBooks struct {
	Featured []FeaturedBook `json:"featured"`
}

type Book struct {
	ID          string     `json:"id"`
	ISBN        string     `json:"isbn"`
	Name        string     `json:"name"`
	Author      string     `json:"author"`
	Type        []BookType `json:"type"`
	Description string     `json:"description"`
	Cover       string     `json:"cover"`
	Genre       string     `json:"genre"`
	Tags        []string   `json:"tags"`
	Link        string     `json:"link"`
	Status      Status     `json:"status"`
	Rating      int        `json:"rating"`
	Pages       int        `json:"pages"`
	Duration    string     `json:"duration"` // For audiobooks
	Publisher   string     `json:"publisher"`
	Published   time.Time  `json:"published"`
	Added       time.Time  `json:"added"`
	Started     time.Time  `json:"started"`
	Finished    time.Time  `json:"finished"`
	Notes       string     `json:"notes"`
	Series      string     `json:"series"`
	SeriesOrder int        `json:"series_order"`
}

type BookType string

const (
	Physical BookType = "physical"
	Audible  BookType = "audible"
	Kindle   BookType = "kindle"
	Ebook    BookType = "ebook"
)

type Status string

const (
	Unread     Status = "unread"
	Reading    Status = "reading"
	Completed  Status = "completed"
	Abandoned  Status = "abandoned"
	WantToRead Status = "want_to_read"
)

type FeaturedBook struct {
	ISBN    string `json:"isbn"`
	Current bool   `json:"current"`
}

type IndexParams struct {
	Host string
}

type PostBook struct {
	Book Book   `json:"book"`
	Key  string `json:"key"`
}

type PostFeatured struct {
	FeaturedBook FeaturedBook `json:"featured_book"`
	Key          string       `json:"key"`
}

type BookFilter struct {
	Type   []BookType `json:"type"`
	Status []Status   `json:"status"`
	Genre  []string   `json:"genre"`
	Author []string   `json:"author"`
	Rating int        `json:"rating"`
	Search string     `json:"search"`
}

type BookStats struct {
	TotalBooks    int            `json:"total_books"`
	ByType        map[BookType]int `json:"by_type"`
	ByStatus      map[Status]int   `json:"by_status"`
	ByGenre       map[string]int   `json:"by_genre"`
	AverageRating float64         `json:"average_rating"`
	PagesRead     int            `json:"pages_read"`
	HoursListened int            `json:"hours_listened"`
}
