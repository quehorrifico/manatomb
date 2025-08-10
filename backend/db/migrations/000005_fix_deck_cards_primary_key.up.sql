-- 000005_fix_deck_cards_primary_key.up.sql

-- First, we drop the old primary key constraint.
ALTER TABLE deck_cards
DROP CONSTRAINT deck_cards_pkey;

-- Then, we add a new composite primary key that includes the 'board'.
-- This allows the same card to exist in a deck on different boards (e.g., main and maybeboard).
ALTER TABLE deck_cards
ADD PRIMARY KEY (deck_id, card_scryfall_id, board);
