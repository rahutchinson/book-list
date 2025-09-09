-- Migration script to convert from SQL table to JSON format
-- This script extracts data from the simple books table and formats it for the JSON-based application

-- Create a temporary table to hold the converted data
CREATE TEMP TABLE books_converted AS
SELECT 
    -- Generate a unique ID for each book (using timestamp + ISBN hash)
    CONCAT(EXTRACT(EPOCH FROM NOW())::TEXT, '_', MD5(isbn)) as id,
    isbn,
    COALESCE(name, 'Unknown Title') as name,
    COALESCE(author, 'Unknown Author') as author,
    -- Map type to the application's expected values
    CASE 
        WHEN LOWER(type) = 'physical' THEN 'physical'
        WHEN LOWER(type) = 'kindle' THEN 'kindle'
        WHEN LOWER(type) = 'audible' THEN 'audible'
        WHEN LOWER(type) = 'ebook' THEN 'ebook'
        ELSE 'physical' -- Default to physical if type is null or unknown
    END as type,
    COALESCE(description, '') as description,
    COALESCE(cover, '') as cover,
    COALESCE(genre, '') as genre,
    -- Convert tags string to array format expected by the application
    CASE 
        WHEN tags IS NOT NULL AND tags != '' THEN 
            ARRAY[tags] -- Single tag as array
        ELSE 
            ARRAY[]::STRING[] -- Empty array
    END as tags,
    COALESCE(link, '') as link,
    'unread' as status, -- Default status
    0 as rating, -- Default rating
    0 as pages, -- Default pages
    '' as duration, -- Default duration
    '' as publisher, -- Default publisher
    NULL as published, -- Default published date
    NOW() as added, -- Current timestamp
    NULL as started, -- Default started date
    NULL as finished, -- Default finished date
    '' as notes, -- Default notes
    '' as series, -- Default series
    0 as series_order -- Default series order
FROM books;

-- Export the converted data as JSON
-- This will create a JSON file that can be used as books.json
SELECT 
    json_build_object(
        'books', 
        json_agg(
            json_build_object(
                'id', id,
                'isbn', isbn,
                'name', name,
                'author', author,
                'type', type,
                'description', description,
                'cover', cover,
                'genre', genre,
                'tags', tags,
                'link', link,
                'status', status,
                'rating', rating,
                'pages', pages,
                'duration', duration,
                'publisher', publisher,
                'published', published,
                'added', added,
                'started', started,
                'finished', finished,
                'notes', notes,
                'series', series,
                'series_order', series_order
            )
        )
    ) as json_output
FROM books_converted;

-- Alternative: Create a CSV export for manual processing
-- COPY (
--     SELECT 
--         id, isbn, name, author, type, description, cover, genre, 
--         array_to_string(tags, ','), link, status, rating, pages, 
--         duration, publisher, published, added, started, finished, 
--         notes, series, series_order
--     FROM books_converted
-- ) TO '/tmp/books_export.csv' WITH CSV HEADER;

-- Print migration summary
SELECT 
    'Migration completed successfully!' as message,
    COUNT(*) as total_books_migrated,
    COUNT(CASE WHEN type = 'physical' THEN 1 END) as physical_books,
    COUNT(CASE WHEN type = 'kindle' THEN 1 END) as kindle_books,
    COUNT(CASE WHEN type = 'audible' THEN 1 END) as audible_books,
    COUNT(CASE WHEN type = 'ebook' THEN 1 END) as ebook_books
FROM books_converted;

-- Clean up
DROP TABLE books_converted;
