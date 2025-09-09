-- Migration script to upgrade from old schema to new enhanced schema
-- Run this script if you have existing data in the old format

-- First, create a backup of existing data
CREATE TABLE books_backup AS SELECT * FROM books;
CREATE TABLE featured_backup AS SELECT * FROM featured;

-- Drop existing tables
DROP TABLE IF EXISTS featured;
DROP TABLE IF EXISTS books;

-- Create new enhanced books table
CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    isbn STRING,
    name STRING NOT NULL,
    author STRING NOT NULL,
    type STRING NOT NULL CHECK (type IN ('physical', 'audible', 'kindle', 'ebook')),
    description STRING,
    cover STRING,
    genre STRING,
    tags STRING[], -- Array of tags
    link STRING,
    status STRING DEFAULT 'unread' CHECK (status IN ('unread', 'reading', 'completed', 'abandoned', 'want_to_read')),
    rating INTEGER CHECK (rating >= 0 AND rating <= 5),
    pages INTEGER,
    duration STRING, -- For audiobooks (e.g., "12h 30m")
    publisher STRING,
    published TIMESTAMP,
    added TIMESTAMP DEFAULT NOW(),
    started TIMESTAMP,
    finished TIMESTAMP,
    notes STRING,
    series STRING,
    series_order INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_books_type ON books(type);
CREATE INDEX idx_books_status ON books(status);
CREATE INDEX idx_books_author ON books(author);
CREATE INDEX idx_books_genre ON books(genre);
CREATE INDEX idx_books_rating ON books(rating);
CREATE INDEX idx_books_series ON books(series);
CREATE INDEX idx_books_added ON books(added);

-- Create featured books table
CREATE TABLE featured (
    isbn STRING PRIMARY KEY,
    current BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Migrate existing data
INSERT INTO books (
    isbn, name, author, type, description, cover, genre, tags, link, status, rating, added
)
SELECT 
    isbn,
    name,
    author,
    CASE 
        WHEN type = 'physical' THEN 'physical'
        WHEN type = 'kindle' THEN 'kindle'
        WHEN type = 'audible' THEN 'audible'
        ELSE 'ebook'
    END as type,
    description,
    cover,
    genre,
    CASE 
        WHEN tags IS NOT NULL AND tags != '' THEN ARRAY[tags]
        ELSE ARRAY[]::STRING[]
    END as tags,
    link,
    'unread' as status,
    0 as rating,
    NOW() as added
FROM books_backup;

-- Migrate featured books
INSERT INTO featured (isbn, current)
SELECT isbn, true FROM featured_backup;

-- Create a view for book statistics
CREATE VIEW book_stats AS
SELECT 
    COUNT(*) as total_books,
    COUNT(CASE WHEN type = 'physical' THEN 1 END) as physical_books,
    COUNT(CASE WHEN type = 'audible' THEN 1 END) as audible_books,
    COUNT(CASE WHEN type = 'kindle' THEN 1 END) as kindle_books,
    COUNT(CASE WHEN type = 'ebook' THEN 1 END) as ebook_books,
    COUNT(CASE WHEN status = 'unread' THEN 1 END) as unread_books,
    COUNT(CASE WHEN status = 'reading' THEN 1 END) as reading_books,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_books,
    COUNT(CASE WHEN status = 'abandoned' THEN 1 END) as abandoned_books,
    COUNT(CASE WHEN status = 'want_to_read' THEN 1 END) as want_to_read_books,
    AVG(rating) as average_rating,
    SUM(CASE WHEN status = 'completed' THEN pages ELSE 0 END) as pages_read
FROM books;

-- Optional: Drop backup tables after confirming migration was successful
-- DROP TABLE books_backup;
-- DROP TABLE featured_backup;

-- Print migration summary
SELECT 
    'Migration completed successfully!' as message,
    COUNT(*) as total_books_migrated
FROM books;
