// src/services/api.js
import axios from 'axios';

const api = axios.create({
  // baseURL: '/api',
  baseURL: process.env.REACT_APP_API_URL,
  headers: { 'Content-Type': 'application/json' },
});

// --- User Auth ---
export const registerUser = (userData) => api.post('/users/register', userData);
export const loginUser = (credentials) => api.post('/users/login', credentials);
export const logoutUser = () => api.post('/users/logout');
export const getCurrentUser = () => api.get('/users/me');

// --- Decks ---
export const createDeck = (deckData) => api.post('/decks/', deckData);
export const getDecks = () => api.get('/decks/');
export const getDeck = (deckId) => api.get(`/decks/${deckId}`);
export const updateDeck = (deckId, deckData) => api.put(`/decks/${deckId}`, deckData);
export const deleteDeck = (deckId) => api.delete(`/decks/${deckId}`);
export const setDeckVisibility = (deckId, isPublic) => api.put(`/decks/${deckId}/visibility`, { is_public: isPublic });

// --- Deck Cards ---
export const addCardToDeck = (deckId, cardData, board) => api.post(`/decks/${deckId}/cards`, { card: cardData, board: board });
export const removeCardFromDeck = (deckId, cardId) => api.delete(`/decks/${deckId}/cards/${cardId}`);


// --- Scryfall API ---
const scryfallApi = axios.create({ baseURL: 'https://api.scryfall.com' });
export const searchScryfall = (query) => scryfallApi.get(`/cards/search?q=${encodeURIComponent(query)}`);

// --- Profiles ---
export const getUserProfile = (username) => api.get(`/profiles/${username}`);

export default api;
