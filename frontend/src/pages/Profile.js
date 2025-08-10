import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getUserProfile } from '../services/api';
import './Profile.css';

function Profile() {
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { username } = useParams();

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const response = await getUserProfile(username);
        setProfile(response.data);
      } catch (err) {
        setError('User profile not found.');
      } finally {
        setLoading(false);
      }
    };
    fetchProfile();
  }, [username]);

  if (loading) return <div>Loading profile...</div>;
  if (error) return <p className="error-message">{error}</p>;
  if (!profile) return null;

  return (
    <div className="profile-container">
      <h1>{profile.username}'s Public Decks</h1>
      <div className="public-decks-list">
        {profile.public_decks.length > 0 ? (
          profile.public_decks.map(deck => (
            <Link to={`/decks/${deck.id}`} key={deck.id} className="deck-link">
              <div className="deck-list-item">
                <h3>{deck.name}</h3>
                <p>Format: {deck.format}</p>
              </div>
            </Link>
          ))
        ) : (
          <p>{profile.username} has not shared any decks yet.</p>
        )}
      </div>
    </div>
  );
}

export default Profile;
