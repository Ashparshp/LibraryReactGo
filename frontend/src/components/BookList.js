import React, { useState } from 'react';
import axios from 'axios';
import './BookList.css';

const BookList = ({ books, fetchBooks }) => {
  const [visibleBooks, setVisibleBooks] = useState(6);
  const [sortOption, setSortOption] = useState('title');

  const borrowBook = async (isbn) => {
    await axios.put(`http://localhost:8080/borrow_book/${isbn}`);
    fetchBooks();
  };

  const returnBook = async (isbn) => {
    await axios.put(`http://localhost:8080/return_book/${isbn}`);
    fetchBooks();
  };

  const removeBook = async (isbn) => {
    await axios.delete(`http://localhost:8080/remove_book/${isbn}`);
    fetchBooks();
  };

  const loadMoreBooks = () => {
    setVisibleBooks(visibleBooks + 6);
  };

  const handleSortChange = (event) => {
    setSortOption(event.target.value);
  };

  const sortBooks = (books, option) => {
    return books.sort((a, b) => {
      if (option === 'title') {
        return a.title.localeCompare(b.title);
      } else if (option === 'author') {
        return a.author.localeCompare(b.author);
      }
      return 0;
    });
  };

  const sortedBooks = sortBooks([...books], sortOption);

  return (
    <div className="book-list-container">
      <h2 className="book-list-heading">Book List</h2>
      <div className="sort-options">
        <label htmlFor="sort">Sort by: </label>
        <select id="sort" value={sortOption} onChange={handleSortChange}>
          <option value="title">Title</option>
          <option value="author">Author</option>
        </select>
      </div>
      {sortedBooks.length > 0 ? (
        sortedBooks.slice(0, visibleBooks).map((book) => (
          <div key={book.isbn} className="book-item">
            <p>
              <span className="book-info">Title:</span> {book.title}<br />
              <span className="book-info">Author:</span> {book.author}<br />
              <span className="book-info">ISBN:</span> {book.isbn}<br />
              <span className="book-info">Quantity:</span> {book.quantity}
            </p>
            <div className="book-buttons">
              <button className="borrow-button" onClick={() => borrowBook(book.isbn)}>Borrow</button>
              <button className="return-button" onClick={() => returnBook(book.isbn)}>Return</button>
              <button className="remove-button" onClick={() => removeBook(book.isbn)}>Remove</button>
            </div>
          </div>
        ))
      ) : (
        <p className="no-books-message">No books available.</p>
      )}
      {sortedBooks.length > visibleBooks && (
        <div className="load-more-container">
          <button className="load-more-button" onClick={loadMoreBooks}>Load More</button>
        </div>
      )}
    </div>
  );
};

export default BookList;
