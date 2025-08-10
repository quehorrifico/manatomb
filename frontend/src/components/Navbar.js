import React from 'react';
import { Link } from 'react-router-dom';
import { logoutUser } from '../services/api';
import './Navbar.css';

function Navbar({ user }) {
  const handleLogout = async () => {
    try {
      await logoutUser();
      window.location.href = '/login';
    } catch (error) {
      console.error('Logout failed', error);
    }
  };

  return (
    <nav className="navbar">
      <Link to="/" className="navbar-brand">Mana Tomb</Link>
      <div className="navbar-links">
        {user ? (
          <>
            <Link to="/my-decks">My Decks</Link>
            <Link to="/search">Card Search</Link>
            <span>Hello, {user.username}</span>
            <button onClick={handleLogout} className="nav-button">Logout</button>
          </>
        ) : (
          <>
            <Link to="/login">Login</Link>
            <Link to="/register">Register</Link>
          </>
        )}
      </div>
    </nav>
  );
}

export default Navbar;
