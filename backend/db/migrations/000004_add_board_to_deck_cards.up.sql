-- 000004_add_board_to_deck_cards.up.sql

-- Adds a 'board' column to the deck_cards table to distinguish
-- between the main deck and other boards like a maybeboard or sideboard.
-- 'main' is the default value for any existing cards.
ALTER TABLE deck_cards
ADD COLUMN board VARCHAR(50) NOT NULL DEFAULT 'main';
