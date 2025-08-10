-- 000002_create_decks_table.up.sql

-- This migration creates the "decks" table.
-- Each deck is associated with a user via the user_id foreign key.
-- ON DELETE CASCADE means if a user is deleted, all their decks are also deleted.
CREATE TABLE IF NOT EXISTS decks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    format VARCHAR(50) NOT NULL DEFAULT 'commander',
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
