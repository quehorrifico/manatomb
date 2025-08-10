import React from 'react';

function Home({ user }) {
  return (
    <div className="home-container">
      {user ? (
        <h1>Welcome to Mana Tomb, {user.username}!</h1>
      ) : (
        <h1>Welcome to Mana Tomb</h1>
      )}
      <p>Your ultimate MTG deck-building companion.</p>
    </div>
  );
}

export default Home;
