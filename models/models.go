package models

type Books struct {
	Books []Book `json:"books"`
}

type FeaturedBooks struct {
	Featured []FeaturedBook `json:"featured"`
}

type Book struct {
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

type FeaturedBook struct {
	ISBN    string `json:"isbn"`
	Current bool   `json:"current"`
}

type IndexParams struct {
	Host string
}

type PostBook struct {
	Book Book
	Key  string
}

type PostFeatured struct {
	FeaturedBook FeaturedBook
	Key          string
}
