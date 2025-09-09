# Virtual Bookshelf

A modern, comprehensive virtual bookshelf application that combines all your physical, Audible, and Kindle books into one beautiful, organized place. Track your reading progress, manage your library, and get insights into your reading habits.

## 🌟 Features

### 📚 Book Management
- **Multi-format support**: Physical books, Kindle, Audible, and E-books
- **Reading status tracking**: Unread, Currently Reading, Completed, Abandoned, Want to Read
- **Rating system**: 5-star rating system for your books
- **Personal notes**: Add your thoughts and notes for each book
- **Cover images**: Beautiful book cover display with 3D effects
- **Series support**: Organize books by series and track reading order

### 🔍 Advanced Filtering & Search
- **Real-time search**: Search by title, author, or description
- **Filter by type**: Physical, Kindle, Audible, E-book
- **Filter by status**: See only books you're reading, completed, etc.
- **Filter by rating**: Find your highest-rated books
- **Genre filtering**: Filter by book genre

### 📊 Reading Statistics
- **Library overview**: Total books, pages read, average rating
- **Format breakdown**: See how many physical vs digital books you have
- **Reading progress**: Track completed vs unread books
- **Visual insights**: Beautiful charts and statistics

### 🎨 Modern UI/UX
- **3D book effects**: Beautiful 3D book covers with hover animations
- **Responsive design**: Works perfectly on desktop, tablet, and mobile
- **Dark/light theme**: Modern gradient backgrounds
- **Smooth animations**: Fluid transitions and interactions
- **Toast notifications**: Real-time feedback for actions

### 🔧 Easy Management
- **Add books**: Simple form to add new books to your library
- **Edit books**: Update book details, status, and ratings
- **Delete books**: Remove books from your library
- **Bulk operations**: Manage multiple books efficiently

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Modern web browser

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/virtual-bookshelf.git
   cd virtual-bookshelf
   ```

2. **Install dependencies and run**
   ```bash
   go mod tidy
   go run main.go
   ```

3. **Optional: Configure environment variables**
   ```bash
   export POST_KEY="your-secret-key-for-api-access"  # Optional
   export PORT="4000"  # Optional, defaults to 4000
   ```

5. **Open your browser**
   Navigate to `http://localhost:4000`

## 📖 Usage Guide

### Adding Books
1. Click the "Add Book" tab in the navigation
2. Fill in the book details (title and author are required)
3. Select the book type (Physical, Kindle, Audible, E-book)
4. Set the reading status
5. Add optional details like genre, pages, cover URL, etc.
6. Click "Add Book" to save

### Managing Your Library
- **Currently Reading**: Books with "Reading" status appear in the featured section
- **Edit books**: Click the edit button on any book to modify details
- **Update status**: Change reading status as you progress through books
- **Add ratings**: Rate books from 1-5 stars
- **Add notes**: Keep personal thoughts about each book

### Using Filters
- **Search**: Type in the search box to find books by title or author
- **Type filter**: Select specific book formats
- **Status filter**: Show only books with certain reading status
- **Rating filter**: Find books with specific ratings
- **Clear filters**: Reset all filters to show all books

### Viewing Statistics
- Click the "Statistics" tab to see your reading insights
- View total books, pages read, average rating
- See breakdown by book type and reading status

## 🏗️ Architecture

### Backend (Go)
- **main.go**: HTTP server and route handlers
- **models/**: Data structures and types

### Frontend (HTML/CSS/JavaScript)
- **index.html**: Main application interface
- **js/main.css**: Modern styling with 3D effects
- **js/main.js**: Interactive functionality and API calls

### Data Storage
- **books.json**: Local JSON file containing all book data
- **Automatic backup**: Data is automatically saved to disk
- **No database required**: Simple file-based storage

## 🔧 Configuration

### Environment Variables
- `POST_KEY`: Secret key for API authentication (optional)
- `PORT`: Server port (default: 4000)

### Data Storage Setup
The application uses local JSON file storage with the following features:
- Unique IDs for books
- Array support for tags
- Automatic data persistence
- Thread-safe file operations
- Built-in data validation

## 🎨 Customization

### Styling
The application uses modern CSS with:
- CSS Grid and Flexbox for layouts
- CSS transforms for 3D effects
- CSS animations for smooth interactions
- Custom properties for theming

### Adding New Features
The modular architecture makes it easy to add new features:
- Add new fields to the Book model
- Create new database migrations
- Add new API endpoints
- Extend the frontend interface

## 📱 Browser Support
- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- 3D book effects inspired by modern CSS techniques
- Icons provided by Font Awesome
- UI framework by Bootstrap
- Database powered by PostgreSQL/CockroachDB

## 🐛 Troubleshooting

### Common Issues

**Data File Error**
- Check that the `books.json` file is writable
- Ensure the application has permission to create/modify files
- Verify the JSON file format is valid

**Books Not Loading**
- Check browser console for JavaScript errors
- Verify the API endpoints are responding
- Ensure the database has data

**Styling Issues**
- Clear browser cache
- Check that all CSS files are loading
- Verify Font Awesome is accessible

### Getting Help
- Check the browser console for error messages
- Verify the `books.json` file exists and is readable
- Check the application logs for backend errors
- Ensure the application has proper file permissions

---

**Happy Reading! 📚✨**