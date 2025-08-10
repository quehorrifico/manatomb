import React, { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import { getDeck, searchScryfall, addCardToDeck, removeCardFromDeck } from '../services/api';
import DeckStats from '../components/DeckStats'; // Import the new component
import './DeckDetail.css';

const CardList = ({ title, cards, onRemove }) => (
  <div className="card-list-section">
    <h2>{title} ({cards.reduce((sum, card) => sum + card.quantity, 0)})</h2>
    {cards.length > 0 ? (
      <div className="card-grid">
        {cards.map((card) => (
          <div key={card.id} className="card-grid-item">
            <img
              src={JSON.parse(card.image_uris)?.normal || ''}
              alt={card.name}
              loading="lazy"
            />
            <div className="card-actions-overlay">
                <button onClick={() => onRemove(card.id)} className="remove-btn">Remove</button>
            </div>
            <span className="quantity-badge">{card.quantity}x</span>
          </div>
        ))}
      </div>
    ) : (
      <p>No cards in this section.</p>
    )}
  </div>
);

const DeckSearch = ({ onAddCard }) => {
    const [query, setQuery] = useState('');
    const [results, setResults] = useState([]);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSearch = async (e) => {
        e.preventDefault();
        if (!query) return;
        setLoading(true);
        setError('');
        try {
            const response = await searchScryfall(query);
            setResults(response.data.data || []);
        } catch (err) {
            setError('No cards found.');
            setResults([]);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="deck-search-section">
            <h2>Add Cards</h2>
            <form onSubmit={handleSearch} className="search-form">
                <input
                    type="text"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    placeholder="Search Scryfall..."
                />
                <button type="submit" disabled={loading}>{loading ? '...' : 'Search'}</button>
            </form>
            {error && <p className="error-message">{error}</p>}
            <div className="search-results-grid">
                {results.map(card => (
                    <div key={card.id} className="search-result-item">
                        <img src={card.image_uris?.small} alt={card.name} loading="lazy" />
                        <div className="search-result-actions">
                            <button onClick={() => onAddCard(card, 'main')}>To Deck</button>
                            <button onClick={() => onAddCard(card, 'maybeboard')}>To Maybe</button>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

function DeckDetail() {
  const [deck, setDeck] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { deckId } = useParams();

  const fetchDeck = useCallback(async () => {
    try {
      const response = await getDeck(deckId);
      setDeck(response.data);
    } catch (err) {
      setError('Failed to fetch deck details.');
    } finally {
      setLoading(false);
    }
  }, [deckId]);

  useEffect(() => {
    fetchDeck();
  }, [fetchDeck]);

  const handleAddCard = async (card, board) => {
    const cardData = {
      id: card.id, name: card.name, image_uris: JSON.stringify(card.image_uris), mana_cost: card.mana_cost,
      cmc: card.cmc, type_line: card.type_line, oracle_text: card.oracle_text, colors: card.colors,
      color_identity: card.color_identity,
    };
    try {
      await addCardToDeck(deckId, cardData, board);
      fetchDeck();
    } catch (err) {
      alert('Failed to add card.');
    }
  };

  const handleRemoveCard = async (cardId) => {
    try {
        await removeCardFromDeck(deckId, cardId);
        fetchDeck();
    } catch (err) {
        alert('Failed to remove card.');
    }
  };

  if (loading) return <div>Loading deck...</div>;
  if (error) return <p className="error-message">{error}</p>;
  if (!deck) return <p>Deck not found.</p>;

  return (
    <div className="deck-detail-container">
      <div className="deck-header">
        <h1>{deck.name}</h1>
        <p>Format: {deck.format}</p>
      </div>

      {/* Render the stats component if there are cards in the mainboard */}
      {deck.mainboard && deck.mainboard.length > 0 && <DeckStats cards={deck.mainboard} />}

      <div className="deck-layout">
        <div className="deck-boards">
            <CardList title="Main Deck" cards={deck.mainboard || []} onRemove={handleRemoveCard} />
            <hr />
            <CardList title="Maybeboard" cards={deck.maybeboard || []} onRemove={handleRemoveCard} />
        </div>
        <div className="deck-sidebar">
            <DeckSearch onAddCard={handleAddCard} />
        </div>
      </div>
    </div>
  );
}

export default DeckDetail;
