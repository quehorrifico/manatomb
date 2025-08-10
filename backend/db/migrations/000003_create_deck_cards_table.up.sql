-- 000003_create_deck_cards_table.up.sql

-- This table links cards to decks, creating a many-to-many relationship.
-- It stores the Scryfall ID for the card and references our local deck ID.
CREATE TABLE IF NOT EXISTS deck_cards (
    deck_id UUID NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
    card_scryfall_id UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    -- A composite primary key ensures a card can only be in a deck once per entry.
    PRIMARY KEY (deck_id, card_scryfall_id)
);

-- We are also adding a new table to cache card data from Scryfall
-- to reduce external API calls and speed up deck loading.
CREATE TABLE IF NOT EXISTS cards (
    scryfall_id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    image_uris JSONB, -- Storing the JSON object of image URIs
    mana_cost VARCHAR(100),
    cmc REAL,
    type_line VARCHAR(255),
    oracle_text TEXT,
    colors VARCHAR(20)[],
    color_identity VARCHAR(20)[]
);
