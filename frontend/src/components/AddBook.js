import React, { useState } from 'react';
import axios from 'axios';
import './AddBook.css';

const AddBook = ({ fetchBooks }) => {
  const [title, setTitle] = useState('');
  const [author, setAuthor] = useState('');
  const [isbn, setIsbn] = useState('');
  const [quantity, setQuantity] = useState('');
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);

  const validateInputs = () => {
    const errors = {};

    if (!title.trim()) {
      errors.title = 'Title is required';
    }

    if (!author.trim()) {
      errors.author = 'Author is required';
    }

    if (!isbn.trim()) {
      errors.isbn = 'ISBN is required';
    } else if (!/^\d{9}$/.test(isbn)) {
      errors.isbn = 'ISBN must be 9 digits';
    }

    if (!quantity) {
      errors.quantity = 'Quantity is required';
    } else if (isNaN(quantity) || parseInt(quantity) <= 0) {
      errors.quantity = 'Quantity must be a positive number';
    }

    setErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const addBook = async () => {
    if (!validateInputs()) return;

    setLoading(true);
    try {
      await axios.post('http://localhost:8080/add_book', { title, author, isbn, quantity: parseInt(quantity) });
      fetchBooks();
      setTitle('');
      setAuthor('');
      setIsbn('');
      setQuantity('');
      setErrors({});
    } catch (err) {
      setErrors({ general: 'Failed to add book. Please try again.' });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="add-book-container">
      <h2 className="add-book-heading">Add Book</h2>
      <div className="add-book-form">
        {errors.general && <p className="error-message">{errors.general}</p>}
        <div className="input-group">
          <input
            type="text"
            placeholder="Title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className={errors.title ? 'input-error' : ''}
          />
          {errors.title && <p className="error-text">{errors.title}</p>}
        </div>
        <div className="input-group">
          <input
            type="text"
            placeholder="Author"
            value={author}
            onChange={(e) => setAuthor(e.target.value)}
            className={errors.author ? 'input-error' : ''}
          />
          {errors.author && <p className="error-text">{errors.author}</p>}
        </div>
        <div className="input-group">
          <input
            type="text"
            placeholder="ISBN"
            value={isbn}
            onChange={(e) => setIsbn(e.target.value.trim())}
            className={errors.isbn ? 'input-error' : ''}
          />
          {errors.isbn && <p className="error-text">{errors.isbn}</p>}
        </div>
        <div className="input-group">
          <input
            type="number"
            placeholder="Quantity"
            value={quantity}
            onChange={(e) => setQuantity(e.target.value)}
            className={errors.quantity ? 'input-error' : ''}
          />
          {errors.quantity && <p className="error-text">{errors.quantity}</p>}
        </div>
        <button className={`add-book-button ${loading ? 'loading' : ''}`} onClick={addBook} disabled={loading}>
          {loading ? 'Adding...' : 'Add Book'}
        </button>
      </div>
    </div>
  );
};

export default AddBook;