import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';
import AddBook from './components/AddBook';
import SearchBooks from './components/SearchBooks';
import BookList from './components/BookList';

function App() {
  const [books, setBooks] = useState([]);
  const [totalBooks, setTotalBooks] = useState(0);

  const fetchBooks = async () => {
    const response = await axios.get('http://localhost:8080/list_books');
    setBooks(response.data);
  };

  const fetchTotalBooks = async () => {
    try {
      const response = await axios.get('http://localhost:8080/total_books');
      setTotalBooks(response.data.totalBooks);
    } catch (error) {
      console.error('Error fetching total books:', error);
    }
  };

  useEffect(() => {
    fetchBooks();
    fetchTotalBooks();
  }, []);

  return (
    <div className="app-container gradient-97">

      <div className="title-container">
        <h1 className="app-title">Library Management System</h1>
        <p className="total-books">Total Books Available: <span className="total-books-number">{totalBooks}</span></p>
      </div>
      <SearchBooks setBooks={setBooks} />
      <BookList books={books} fetchBooks={fetchBooks} />
      <AddBook fetchBooks={fetchBooks} />
    </div>
  );
}

export default App;
