import React, { useState } from 'react';
import axios from 'axios';
import './SearchBooks.css';

const SearchBooks = ({ setBooks }) => {
  const [query, setQuery] = useState('');
  const [searchResult, setSearchResult] = useState(null); // Use null instead of an empty array

  const searchBooks = async () => {
    if (query.trim() === '') {
      setSearchResult(null); // Clear the search results if the query is empty
      return;
    }

    try {
      const response = await axios.get('http://localhost:8080/search_book', { params: { query } });
      setBooks(response.data); // Update parent component state with search results
      setSearchResult(response.data); // Update local state for display purposes
    } catch (error) {
      console.error('Error searching books:', error);
      // Handle error (e.g., show error message)
    }
  };

  const borrowBook = async (isbn) => {
    try {
      await axios.put(`http://localhost:8080/borrow_book/${isbn}`);
      // After borrowing, refresh the search results
      searchBooks();
    } catch (error) {
      console.error('Error borrowing book:', error);
      // Handle error (e.g., show error message)
    }
  };

  return (
    <div className="search-books-container">
      <h2 className="search-books-heading">Search Books</h2>
      <div className="search-books-form">
        <input
          type="text"
          className="search-books-input"
          placeholder="Search by Title, Author or ISBN"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <button className="search-books-button" onClick={searchBooks}>Search</button>
      </div>

      {/* Display search results */}
      {searchResult !== null && ( // Check if searchResult is not null
        <div className="search-results">
          <h3 className="search-results-heading">Search Results</h3>
          {searchResult.map((book) => (
            <div key={book.isbn} className="search-result-item">
              <p><span className="search-result-label">Title:</span> {book.title}</p>
              <p><span className="search-result-label">Author:</span> {book.author}</p>
              <p><span className="search-result-label">ISBN:</span> {book.isbn}</p>
              <p><span className="search-result-label">Quantity:</span> {book.quantity}</p>
              <button className="borrow-button" onClick={() => borrowBook(book.isbn)}>Borrow</button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default SearchBooks;
