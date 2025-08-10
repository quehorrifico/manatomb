-- 000006_add_public_to_decks.up.sql

-- Adds a boolean column to mark decks as public or private.
-- It defaults to FALSE, so all existing and new decks are private by default.
ALTER TABLE decks
ADD COLUMN is_public BOOLEAN NOT NULL DEFAULT FALSE;
