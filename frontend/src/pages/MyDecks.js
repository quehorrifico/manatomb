import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { getDecks, createDeck, updateDeck, deleteDeck, setDeckVisibility } from '../services/api';
import './MyDecks.css';

function MyDecks() {
  const [decks, setDecks] = useState([]);
  const [deckName, setDeckName] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);

  const [editingDeckId, setEditingDeckId] = useState(null);
  const [editedName, setEditedName] = useState('');

  const fetchDecks = async () => {
    try {
      const response = await getDecks();
      // The response now includes the is_public flag for each deck
      setDecks(response.data);
    } catch (err) {
      setError('Failed to fetch decks.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDecks();
  }, []);

  const handleCreateDeck = async (e) => {
    e.preventDefault();
    if (!deckName) {
      setError('Deck name is required.');
      return;
    }
    setError('');
    try {
      await createDeck({ name: deckName, format: 'commander' });
      setDeckName('');
      fetchDecks();
    } catch (err) {
      setError('Failed to create deck.');
    }
  };

  const handleDeleteDeck = async (deckId) => {
    if (window.confirm('Are you sure you want to delete this deck?')) {
        try {
            await deleteDeck(deckId);
            fetchDecks();
        } catch (err) {
            setError('Failed to delete deck.');
        }
    }
  };

  const handleEditClick = (deck) => {
    setEditingDeckId(deck.id);
    setEditedName(deck.name);
  };

  const handleCancelEdit = () => {
    setEditingDeckId(null);
    setEditedName('');
  };

  const handleSaveEdit = async (deckId) => {
    try {
        await updateDeck(deckId, { name: editedName });
        setEditingDeckId(null);
        fetchDecks();
    } catch (err) {
        setError('Failed to update deck.');
    }
  }

  // This is the function that was missing
  const handleVisibilityToggle = async (deckId, currentStatus) => {
    try {
        await setDeckVisibility(deckId, !currentStatus);
        // Update the state locally for an instant UI change
        setDecks(decks.map(d => 
            d.id === deckId ? { ...d, is_public: !currentStatus } : d
        ));
    } catch (err) {
        setError('Failed to update deck visibility.');
        // If the API call fails, revert the change in the UI
        setDecks(decks.map(d => 
            d.id === deckId ? { ...d, is_public: currentStatus } : d
        ));
    }
  };

  if (loading) return <div>Loading decks...</div>;

  return (
    <div className="my-decks-container">
      <div className="create-deck-form">
        <h2>Create a New Deck</h2>
        <form onSubmit={handleCreateDeck}>
          <input
            type="text"
            value={deckName}
            onChange={(e) => setDeckName(e.target.value)}
            placeholder="Enter deck name"
          />
          <button type="submit">Create Deck</button>
        </form>
        {error && <p className="error-message">{error}</p>}
      </div>

      <div className="decks-list">
        <h2>My Decks</h2>
        {decks.length > 0 ? (
          <ul>
            {decks.map((deck) => (
              <li key={deck.id} className="deck-list-item">
                {editingDeckId === deck.id ? (
                  <div className="edit-deck-form">
                    <input 
                      type="text" 
                      value={editedName} 
                      onChange={(e) => setEditedName(e.target.value)}
                    />
                    <button onClick={() => handleSaveEdit(deck.id)} className="save-btn">Save</button>
                    <button onClick={handleCancelEdit} className="cancel-btn">Cancel</button>
                  </div>
                ) : (
                  <>
                    <Link to={`/decks/${deck.id}`} className="deck-link-content">
                      <h3>{deck.name}</h3>
                      <p>Format: {deck.format}</p>
                    </Link>
                    <div className="deck-actions">
                      <div className="visibility-toggle">
                          <span>Public</span>
                          <label className="switch">
                              <input 
                                  type="checkbox" 
                                  checked={deck.is_public} 
                                  onChange={() => handleVisibilityToggle(deck.id, deck.is_public)}
                              />
                              <span className="slider round"></span>
                          </label>
                      </div>
                      <button onClick={() => handleEditClick(deck)} className="edit-btn">Edit</button>
                      <button onClick={() => handleDeleteDeck(deck.id)} className="delete-btn">Delete</button>
                    </div>
                  </>
                )}
              </li>
            ))}
          </ul>
        ) : (
          <p>You don't have any decks yet. Create one above!</p>
        )}
      </div>
    </div>
  );
}

export default MyDecks;
