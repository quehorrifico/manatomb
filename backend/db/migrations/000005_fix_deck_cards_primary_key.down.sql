-- 000005_fix_deck_cards_primary_key.down.sql

-- This reverts the changes from the .up.sql file.
-- We drop the new composite primary key.
ALTER TABLE deck_cards
DROP CONSTRAINT deck_cards_pkey;

-- And we restore the old primary key.
ALTER TABLE deck_cards
ADD PRIMARY KEY (deck_id, card_scryfall_id);
