# Migration Guide: SQL Table to JSON Format

This guide explains how to migrate data from the simple SQL table structure to the JSON format used by the book-list application.

## SQL Table Structure

The source SQL table has the following structure:

```sql
CREATE TABLE public.books (
  isbn STRING NOT NULL,
  name STRING NULL,
  author STRING NULL,
  type STRING NULL,
  description STRING NULL,
  cover STRING NULL,
  genre STRING NULL,
  tags STRING NULL,
  link STRING NULL,
  CONSTRAINT "primary" PRIMARY KEY (isbn ASC)
)
```

## Migration Options

### Option 1: SQL Migration Script (`migration_from_sql.sql`)

This script converts the SQL table data to JSON format directly in the database.

**Usage:**
1. Connect to your PostgreSQL database
2. Run the migration script:
   ```bash
   psql -d your_database -f migration_from_sql.sql
   ```
3. The script will output JSON that you can save as `books.json`

**Features:**
- Generates unique IDs for each book
- Maps book types to application enum values
- Converts tags string to array format
- Sets default values for missing fields
- Provides migration statistics

### Option 2: Go Migration Script (`migrate.go`)

This Go script connects to CockroachDB and generates the JSON file directly.

**Prerequisites:**
```bash
go get github.com/lib/pq
```

**Usage:**
1. Set environment variables for database connection (optional - defaults are provided):
   ```bash
   export DB_HOST=loyal-efreet-5669.5xj.gcp-us-central1.cockroachlabs.cloud
   export DB_PORT=26257
   export DB_USER=ryan
   export DB_PASSWORD=vBbmP0zJJQsz5Q_0Qzdymw
   export DB_NAME=defaultdb
   export OUTPUT_FILE=books.json
   ```

2. Run the migration:
   ```bash
   go run migrate.go
   ```

**Features:**
- Direct database connection
- Automatic JSON file generation
- Error handling and logging
- Configurable via environment variables

## Field Mapping

| SQL Field | JSON Field | Notes |
|-----------|------------|-------|
| `isbn` | `isbn` | Direct mapping |
| `name` | `name` | Direct mapping |
| `author` | `author` | Direct mapping |
| `type` | `type` | Mapped to enum: physical, kindle, audible, ebook |
| `description` | `description` | Direct mapping |
| `cover` | `cover` | Direct mapping |
| `genre` | `genre` | Direct mapping |
| `tags` | `tags` | Converted from string to array |
| `link` | `link` | Direct mapping |

## Default Values

For fields not present in the source SQL table, the following defaults are applied:

- `id`: Generated unique ID (timestamp + ISBN hash)
- `status`: "unread"
- `rating`: 0
- `pages`: 0
- `duration`: ""
- `publisher`: ""
- `published`: null
- `added`: Current timestamp
- `started`: null
- `finished`: null
- `notes`: ""
- `series`: ""
- `series_order`: 0

## Type Mapping

The SQL `type` field is mapped to the application's BookType enum:

- `physical` → `physical`
- `kindle` → `kindle`
- `audible` → `audible`
- `ebook` → `ebook`
- Any other value → `physical` (default)

## Tags Conversion

The SQL `tags` field (string) is converted to the JSON `tags` field (array):
- If tags is not null/empty: `["tag_value"]`
- If tags is null/empty: `[]`

## Verification

After migration, verify the data by:

1. Checking the generated `books.json` file
2. Running the application and viewing the books
3. Comparing the count of migrated books with the source table

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Verify database credentials
   - Check if PostgreSQL is running
   - Ensure the database exists

2. **Permission Denied**
   - Check file write permissions for the output directory
   - Ensure database user has SELECT permissions on the books table

3. **Invalid JSON Output**
   - Check for special characters in text fields
   - Verify the source data doesn't contain invalid UTF-8 sequences

### Debug Mode

For the Go script, you can add debug logging by modifying the script to print each book as it's processed.

## Rollback

To rollback the migration:
1. Keep a backup of your original `books.json` file
2. Restore from backup if needed
3. The migration scripts don't modify the source SQL table, so your original data is safe

## Support

If you encounter issues with the migration:
1. Check the database connection parameters
2. Verify the source table structure matches the expected format
3. Review the migration logs for specific error messages
