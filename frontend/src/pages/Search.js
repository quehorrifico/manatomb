import React, { useState, useEffect } from 'react';
import { searchScryfall, getDecks, addCardToDeck } from '../services/api';
import './Search.css';

function Search() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [decks, setDecks] = useState([]);
  const [selectedDeck, setSelectedDeck] = useState('');
  const [addCardMessage, setAddCardMessage] = useState('');

  // Fetch user's decks when the component mounts
  useEffect(() => {
    const fetchDecks = async () => {
      try {
        const response = await getDecks();
        setDecks(response.data);
        if (response.data.length > 0) {
          setSelectedDeck(response.data[0].id);
        }
      } catch (err) {
        console.error('Failed to fetch decks for search page');
      }
    };
    fetchDecks();
  }, []);

  const handleSearch = async (e) => {
    e.preventDefault();
    if (!query) return;
    setLoading(true);
    setError('');
    setResults([]);
    setAddCardMessage('');
    try {
      const response = await searchScryfall(query);
      if (response.data && response.data.data) {
        setResults(response.data.data);
      } else {
        setError('No cards found for your query.');
      }
    } catch (err) {
      setError('Failed to fetch card data. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleAddCard = async (card) => {
    if (!selectedDeck) {
      setAddCardMessage('Please select a deck first.');
      return;
    }
    setAddCardMessage(`Adding ${card.name}...`);

    // We need to structure the card data to match our backend model
    const cardData = {
      id: card.id,
      name: card.name,
      image_uris: JSON.stringify(card.image_uris),
      mana_cost: card.mana_cost,
      cmc: card.cmc,
      type_line: card.type_line,
      oracle_text: card.oracle_text,
      colors: card.colors,
      color_identity: card.color_identity,
    };

    try {
      await addCardToDeck(selectedDeck, cardData);
      setAddCardMessage(`Successfully added ${card.name} to your deck!`);
    } catch (err) {
      setAddCardMessage(`Failed to add ${card.name}.`);
    }
  };

  return (
    <div className="search-container">
      <h1>Card Search</h1>
      <form onSubmit={handleSearch} className="search-form">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search for a card (e.g., Sol Ring)"
        />
        <button type="submit" disabled={loading}>
          {loading ? 'Searching...' : 'Search'}
        </button>
      </form>

      {addCardMessage && <p className="add-card-message">{addCardMessage}</p>}
      {error && <p className="error-message">{error}</p>}

      <div className="search-results">
        {results.map((card) => (
          <div key={card.id} className="card-item">
            <img
              src={card.image_uris?.normal || 'https://placehold.co/223x310/1a1a1a/e0e0e0?text=No+Image'}
              alt={card.name}
              loading="lazy"
            />
            <div className="add-to-deck-container">
              <select value={selectedDeck} onChange={(e) => setSelectedDeck(e.target.value)} disabled={decks.length === 0}>
                {decks.length > 0 ? (
                  decks.map(deck => <option key={deck.id} value={deck.id}>{deck.name}</option>)
                ) : (
                  <option>No decks available</option>
                )}
              </select>
              <button onClick={() => handleAddCard(card)} disabled={decks.length === 0}>
                Add
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default Search;
