$(document).ready(function() {
    function getBooks() {
        $.getJSON('/books', function(data) {
            var list = document.querySelector('#bookList');
            $.each(data.books, function(key, value) {
            list.appendChild(buildBook(value))
        });
            
        });
    }
    getBooks();
});



function buildBook(value) {
    var template = document.querySelector('#bookrow');

    var clone = template.content.cloneNode(true);
    var title = clone.querySelector("#title")
    var author = clone.querySelector("#author")
    var cover = clone.querySelector("#cover")
    var link = clone.querySelector("#link")
    var bookType = clone.querySelector("#type")
    title.textContent = value.name;
    author.textContent = value.author;
    cover.src = value.cover
    link.href = value.link
    bookType.textContent = value.type
    return clone;
}