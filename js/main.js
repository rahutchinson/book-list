$(document).ready(function() {
    let allBooks = [];
    let currentView = 'bookshelf';
    
    // Initialize the application
    init();
    
    function init() {
        loadBooks();
        setupEventListeners();
        setupNavigation();
    }
    
    function setupEventListeners() {
        // Navigation
        $('.nav-link').on('click', function(e) {
            e.preventDefault();
            const view = $(this).data('view');
            switchView(view);
        });
        
        // Refresh button
        $('#refreshBtn').on('click', function() {
            loadBooks();
            showToast('Books refreshed!', 'success');
        });
        
        // Add book form
        $('#addBookForm').on('submit', function(e) {
            e.preventDefault();
            addBook();
        });
        
        // Filter controls
        $('#applyFilters').on('click', function() {
            applyFilters();
        });
        
        $('#clearFilters').on('click', function() {
            clearFilters();
        });
        
        // Search input
        $('#searchInput').on('input', function() {
            debounce(applyFilters, 300)();
        });
        
        // Edit book form
        $('#editBookForm').on('submit', function(e) {
            e.preventDefault();
            saveBookChanges();
        });
        
        $('#resetBookChanges').on('click', function() {
            resetBookChanges();
        });
        
        // Cover preview in edit modal
        $('#editBookCover').on('input', function() {
            const coverUrl = $(this).val();
            const preview = $('#editBookCoverPreview');
            
            if (coverUrl && coverUrl.trim()) {
                preview.attr('src', coverUrl).show();
            } else {
                preview.hide();
            }
        });
        
        // ISBN lookup for add book form
        $('#lookupISBN').on('click', function() {
            const isbn = $('#bookISBN').val().trim();
            if (isbn) {
                lookupBookByISBN(isbn, 'add');
            } else {
                showToast('Please enter an ISBN', 'error');
            }
        });
        
        // ISBN lookup for edit book form
        $('#editLookupISBN').on('click', function() {
            const isbn = $('#editBookISBN').val().trim();
            if (isbn) {
                lookupBookByISBN(isbn, 'edit');
            } else {
                showToast('Please enter an ISBN', 'error');
            }
        });
        
        // Handle modal close with unsaved changes
        $('#editBookModal').on('hide.bs.modal', function(e) {
            const originalData = $('#editBookModal').data('original-book');
            if (!originalData) return;
            
            const currentData = {
                name: $('#editBookTitle').val(),
                author: $('#editBookAuthor').val(),
                type: $('#editBookType').val(),
                status: $('#editBookStatus').val(),
                rating: parseInt($('#editBookRating').val()) || 0,
                genre: $('#editBookGenre').val(),
                pages: parseInt($('#editBookPages').val()) || 0,
                cover: $('#editBookCover').val(),
                link: $('#editBookLink').val(),
                description: $('#editBookDescription').val(),
                notes: $('#editBookNotes').val()
            };
            
            // Check if data has changed
            const hasChanges = JSON.stringify(originalData) !== JSON.stringify(currentData);
            
            if (hasChanges) {
                if (!confirm('You have unsaved changes. Are you sure you want to close?')) {
                    e.preventDefault();
                    return;
                }
            }
            
            // Clear the original data and reset form
            $('#editBookModal').removeData('original-book');
            $('#editBookForm')[0].reset();
            $('#editBookCoverPreview').hide();
        });
    }
    
    function setupNavigation() {
        $('.nav-link').removeClass('active');
        $(`.nav-link[data-view="${currentView}"]`).addClass('active');
    }
    
    function switchView(view) {
        $('.view-section').hide();
        $(`#${view}View`).show();
        currentView = view;
        setupNavigation();
        
        if (view === 'stats') {
            loadStats();
        }
    }
    
    function loadBooks() {
        showLoading('#bookList');
        showLoading('#currentlyReading');
        
        $.getJSON('/books', function(data) {
            allBooks = data.books || [];
            renderBooks();
        }).fail(function() {
            showError('Failed to load books');
        });
    }
    
    function renderBooks() {
        const currentlyReading = allBooks.filter(book => book.status === 'reading');
        const otherBooks = allBooks.filter(book => book.status !== 'reading');
        
        renderBookSection('#currentlyReading', currentlyReading);
        renderBookSection('#bookList', otherBooks);
        
        if (currentlyReading.length === 0) {
            $('#currentlyReading').html('<div class="empty-state"><i class="fas fa-book-open"></i><h3>No books currently being read</h3><p>Start reading a book to see it here!</p></div>');
        }
        
        if (otherBooks.length === 0) {
            $('#bookList').html('<div class="empty-state"><i class="fas fa-books"></i><h3>No books in your library</h3><p>Add your first book to get started!</p></div>');
        }
    }
    
    function renderBookSection(container, books) {
        const $container = $(container);
        $container.empty();
        
        if (books.length === 0) return;
        
        // Create the shelf structure
        const $shelf = $('<div class="bookshelf-shelf"></div>');
        $container.append($shelf);
        
        books.forEach(book => {
            const bookElement = createBookElement(book);
            $shelf.append(bookElement);
        });
        
        // Setup book action listeners
        setupBookActions(container);
    }
    
    function createBookElement(book) {
        const template = document.querySelector('#bookTemplate');
        const clone = template.content.cloneNode(true);
        
        // Set book cover
        const coverImg = clone.querySelector('.book-cover');
        coverImg.src = book.cover || 'https://via.placeholder.com/160x240/f8f9fa/6c757d?text=No+Cover';
        coverImg.alt = book.name;
        
        // Set book info
        clone.querySelector('.book-title').textContent = book.name;
        clone.querySelector('.book-author').textContent = book.author;
        clone.querySelector('.book-genre').textContent = book.genre || 'No genre';
        
        // Set book type
        const typeElement = clone.querySelector('.book-type');
        typeElement.textContent = getTypeDisplayName(book.type);
        // Handle multiple types for CSS classes
        if (Array.isArray(book.type)) {
            typeElement.className = `book-type ${book.type.join(' ')}`;
        } else {
            typeElement.className = `book-type ${book.type}`;
        }
        
        // Set book status
        const statusElement = clone.querySelector('.book-status');
        statusElement.textContent = getStatusDisplayName(book.status);
        statusElement.className = `book-status ${book.status}`;
        
        // Set rating
        const ratingElement = clone.querySelector('.book-rating');
        ratingElement.innerHTML = createStarRating(book.rating);
        
        // Store book data
        const bookItem = clone.querySelector('.book-item');
        bookItem.dataset.bookId = book.id;
        bookItem.dataset.bookData = JSON.stringify(book);
        
        return clone;
    }
    
    function setupBookActions(container) {
        $(container).find('.edit-book').on('click', function() {
            const bookItem = $(this).closest('.book-item');
            const bookDataString = bookItem.attr('data-book-data');
            
            if (!bookDataString) {
                console.error('No book data found for edit button');
                showToast('Error: Book data not found', 'error');
                return;
            }
            
            try {
                const bookData = JSON.parse(bookDataString);
                openEditModal(bookData);
            } catch (error) {
                console.error('Failed to parse book data:', error, 'Data:', bookDataString);
                showToast('Error: Invalid book data', 'error');
            }
        });
        
        $(container).find('.delete-book').on('click', function() {
            const bookItem = $(this).closest('.book-item');
            const bookDataString = bookItem.attr('data-book-data');
            
            if (!bookDataString) {
                console.error('No book data found for delete button');
                showToast('Error: Book data not found', 'error');
                return;
            }
            
            try {
                const bookData = JSON.parse(bookDataString);
                deleteBook(bookData);
            } catch (error) {
                console.error('Failed to parse book data:', error, 'Data:', bookDataString);
                showToast('Error: Invalid book data', 'error');
            }
        });
    }
    
    function getTypeDisplayName(type) {
        const typeNames = {
            'physical': 'Physical',
            'kindle': 'Kindle',
            'audible': 'Audible',
            'ebook': 'E-Book'
        };
        
        if (Array.isArray(type)) {
            return type.map(t => typeNames[t] || t).join(', ');
        }
        return typeNames[type] || type;
    }
    
    function getStatusDisplayName(status) {
        const statusNames = {
            'unread': 'Unread',
            'reading': 'Reading',
            'completed': 'Completed',
            'abandoned': 'Abandoned',
            'want_to_read': 'Want to Read'
        };
        return statusNames[status] || status;
    }
    
    function createStarRating(rating) {
        if (!rating || rating === 0) return '<span class="text-muted">No rating</span>';
        
        let stars = '';
        for (let i = 1; i <= 5; i++) {
            if (i <= rating) {
                stars += '<i class="fas fa-star star"></i>';
            } else {
                stars += '<i class="far fa-star star"></i>';
            }
        }
        return stars;
    }
    
    function addBook() {
        const bookData = {
            name: $('#bookTitle').val(),
            author: $('#bookAuthor').val(),
            isbn: $('#bookISBN').val(),
            type: $('#bookType').val() || [],
            status: $('#bookStatus').val(),
            rating: parseInt($('#bookRating').val()) || 0,
            genre: $('#bookGenre').val(),
            pages: parseInt($('#bookPages').val()) || 0,
            cover: $('#bookCover').val(),
            link: $('#bookLink').val(),
            description: $('#bookDescription').val(),
            notes: $('#bookNotes').val(),
            added: new Date().toISOString()
        };
        
        $.ajax({
            url: '/books',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                book: bookData,
                key: '' // No key required for local JSON storage
            }),
            success: function() {
                showToast('Book added successfully!', 'success');
                $('#addBookForm')[0].reset();
                loadBooks();
            },
            error: function() {
                showToast('Failed to add book', 'error');
            }
        });
    }
    
    function openEditModal(book) {
        $('#editBookId').val(book.id);
        $('#editBookTitle').val(book.name);
        $('#editBookAuthor').val(book.author);
        $('#editBookISBN').val(book.isbn || '');
        // Handle array of types
        if (Array.isArray(book.type)) {
            $('#editBookType').val(book.type);
        } else {
            $('#editBookType').val([book.type]);
        }
        $('#editBookStatus').val(book.status);
        $('#editBookRating').val(book.rating || 0);
        $('#editBookGenre').val(book.genre || '');
        $('#editBookPages').val(book.pages || '');
        $('#editBookCover').val(book.cover || '');
        $('#editBookLink').val(book.link || '');
        $('#editBookDescription').val(book.description || '');
        $('#editBookNotes').val(book.notes || '');
        
        // Show cover preview if available
        const preview = $('#editBookCoverPreview');
        if (book.cover && book.cover.trim()) {
            preview.attr('src', book.cover).show();
        } else {
            preview.hide();
        }
        
        // Store original book data for change detection
        $('#editBookModal').data('original-book', {
            name: book.name,
            author: book.author,
            isbn: book.isbn || '',
            type: book.type,
            status: book.status,
            rating: book.rating || 0,
            genre: book.genre || '',
            pages: book.pages || 0,
            cover: book.cover || '',
            link: book.link || '',
            description: book.description || '',
            notes: book.notes || ''
        });
        
        $('#editBookModal').modal('show');
    }
    
    function saveBookChanges() {
        const bookData = {
            id: $('#editBookId').val(),
            name: $('#editBookTitle').val(),
            author: $('#editBookAuthor').val(),
            isbn: $('#editBookISBN').val(),
            type: $('#editBookType').val() || [],
            status: $('#editBookStatus').val(),
            rating: parseInt($('#editBookRating').val()) || 0,
            genre: $('#editBookGenre').val(),
            pages: parseInt($('#editBookPages').val()) || 0,
            cover: $('#editBookCover').val(),
            link: $('#editBookLink').val(),
            description: $('#editBookDescription').val(),
            notes: $('#editBookNotes').val()
        };
        
        // Validate required fields
        if (!bookData.name.trim()) {
            showToast('Book title is required', 'error');
            return;
        }
        if (!bookData.author.trim()) {
            showToast('Book author is required', 'error');
            return;
        }
        
        $.ajax({
            url: '/books',
            method: 'PUT',
            contentType: 'application/json',
            data: JSON.stringify({
                book: bookData,
                key: '' // No key required for local JSON storage
            }),
            success: function() {
                showToast('Book updated successfully!', 'success');
                $('#editBookModal').modal('hide');
                $('#editBookForm')[0].reset();
                $('#editBookModal').removeData('original-book');
                loadBooks();
            },
            error: function() {
                showToast('Failed to update book', 'error');
            }
        });
    }
    
    function resetBookChanges() {
        const originalData = $('#editBookModal').data('original-book');
        if (!originalData) return;
        
        $('#editBookTitle').val(originalData.name);
        $('#editBookAuthor').val(originalData.author);
        $('#editBookISBN').val(originalData.isbn || '');
        // Handle array of types
        if (Array.isArray(originalData.type)) {
            $('#editBookType').val(originalData.type);
        } else {
            $('#editBookType').val([originalData.type]);
        }
        $('#editBookStatus').val(originalData.status);
        $('#editBookRating').val(originalData.rating);
        $('#editBookGenre').val(originalData.genre);
        $('#editBookPages').val(originalData.pages);
        $('#editBookCover').val(originalData.cover);
        $('#editBookLink').val(originalData.link);
        $('#editBookDescription').val(originalData.description);
        $('#editBookNotes').val(originalData.notes);
        
        // Update cover preview
        const preview = $('#editBookCoverPreview');
        if (originalData.cover && originalData.cover.trim()) {
            preview.attr('src', originalData.cover).show();
        } else {
            preview.hide();
        }
        
        showToast('Changes reset to original values', 'info');
    }
    
    function deleteBook(book) {
        if (!confirm(`Are you sure you want to delete "${book.name}"?`)) {
            return;
        }
        
        $.ajax({
            url: '/books',
            method: 'DELETE',
            contentType: 'application/json',
            data: JSON.stringify({
                book: book,
                key: '' // No key required for local JSON storage
            }),
            success: function() {
                showToast('Book deleted successfully!', 'success');
                loadBooks();
            },
            error: function() {
                showToast('Failed to delete book', 'error');
            }
        });
    }
    
    function applyFilters() {
        const filter = {
            search: $('#searchInput').val(),
            type: $('#typeFilter').val() ? [$('#typeFilter').val()] : [],
            status: $('#statusFilter').val() ? [$('#statusFilter').val()] : [],
            rating: parseInt($('#ratingFilter').val()) || 0
        };
        
        $.ajax({
            url: '/books/filter',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(filter),
            success: function(data) {
                allBooks = data.books || [];
                renderBooks();
            },
            error: function() {
                showToast('Failed to apply filters', 'error');
            }
        });
    }
    
    function clearFilters() {
        $('#searchInput').val('');
        $('#typeFilter').val('');
        $('#statusFilter').val('');
        $('#ratingFilter').val('0');
        loadBooks();
    }
    
    function loadStats() {
        $.getJSON('/books/stats', function(data) {
            renderStats(data);
        }).fail(function() {
            showError('Failed to load statistics');
        });
    }
    
    function renderStats(stats) {
        const statsHtml = `
            <div class="col-md-3 mb-4">
                <div class="stats-card">
                    <div class="stats-number">${stats.total_books}</div>
                    <div class="stats-label">Total Books</div>
                </div>
            </div>
            <div class="col-md-3 mb-4">
                <div class="stats-card">
                    <div class="stats-number">${stats.physical_books}</div>
                    <div class="stats-label">Physical Books</div>
                </div>
            </div>
            <div class="col-md-3 mb-4">
                <div class="stats-card">
                    <div class="stats-number">${stats.audible_books}</div>
                    <div class="stats-label">Audiobooks</div>
                </div>
            </div>
            <div class="col-md-3 mb-4">
                <div class="stats-card">
                    <div class="stats-number">${stats.kindle_books + stats.ebook_books}</div>
                    <div class="stats-label">E-Books</div>
                </div>
            </div>
            <div class="col-md-3 mb-4">
                <div class="stats-card">
                    <div class="stats-number">${stats.completed_books}</div>
                    <div class="stats-label">Completed</div>
                </div>
            </div>
            <div class="col-md-3 mb-4">
                <div class="stats-card">
                    <div class="stats-number">${stats.reading_books}</div>
                    <div class="stats-label">Currently Reading</div>
                </div>
            </div>
            <div class="col-md-3 mb-4">
                <div class="stats-card">
                    <div class="stats-number">${stats.pages_read.toLocaleString()}</div>
                    <div class="stats-label">Pages Read</div>
                </div>
            </div>
            <div class="col-md-3 mb-4">
                <div class="stats-card">
                    <div class="stats-number">${stats.average_rating ? stats.average_rating.toFixed(1) : 'N/A'}</div>
                    <div class="stats-label">Avg Rating</div>
                </div>
            </div>
        `;
        
        $('#statsContent').html(statsHtml);
    }
    
    function showLoading(container) {
        $(container).html('<div class="loading"><div class="spinner"></div></div>');
    }
    
    function showError(message) {
        showToast(message, 'error');
    }
    
    function showToast(message, type = 'info') {
        const toastClass = type === 'error' ? 'bg-danger' : type === 'success' ? 'bg-success' : 'bg-info';
        const toast = $(`
            <div class="toast show ${toastClass} text-white" role="alert">
                <div class="toast-body">
                    ${message}
                </div>
            </div>
        `);
        
        $('body').append(toast);
        
        setTimeout(() => {
            toast.remove();
        }, 3000);
    }
    
    function debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }
    
    function lookupBookByISBN(isbn, formType) {
        // Show loading state
        const button = formType === 'add' ? $('#lookupISBN') : $('#editLookupISBN');
        const originalText = button.html();
        button.html('<i class="fas fa-spinner fa-spin"></i> Looking up...').prop('disabled', true);
        
        $.ajax({
            url: '/books/lookup',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ isbn: isbn }),
            success: function(data) {
                if (data.success && data.book) {
                    // Populate form fields with book data
                    if (formType === 'add') {
                        $('#bookTitle').val(data.book.title || '');
                        $('#bookAuthor').val(data.book.author || '');
                        $('#bookGenre').val(data.book.genre || '');
                        $('#bookPages').val(data.book.pages || '');
                        $('#bookCover').val(data.book.cover || '');
                        $('#bookDescription').val(data.book.description || '');
                    } else {
                        $('#editBookTitle').val(data.book.title || '');
                        $('#editBookAuthor').val(data.book.author || '');
                        $('#editBookGenre').val(data.book.genre || '');
                        $('#editBookPages').val(data.book.pages || '');
                        $('#editBookCover').val(data.book.cover || '');
                        $('#editBookDescription').val(data.book.description || '');
                        
                        // Update cover preview
                        const preview = $('#editBookCoverPreview');
                        if (data.book.cover && data.book.cover.trim()) {
                            preview.attr('src', data.book.cover).show();
                        } else {
                            preview.hide();
                        }
                    }
                    
                    showToast('Book details populated successfully!', 'success');
                } else {
                    showToast(data.message || 'Book not found', 'error');
                }
            },
            error: function() {
                showToast('Failed to lookup book details', 'error');
            },
            complete: function() {
                // Restore button state
                button.html(originalText).prop('disabled', false);
            }
        });
    }
});
