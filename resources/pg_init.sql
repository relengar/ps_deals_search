CREATE EXTENSION vector;

-- The embedding model we use (all-MiniLM-L6-v2) uses 384 dimensions
CREATE TABLE games (
    id bigserial PRIMARY KEY,
    name text,
    description text,
    price real,
    original_price real,
    url text,
    rating real,
    rating_sum int,
    expiration timestamp
);

-- The embedding model we use (all-MiniLM-L6-v2) uses 384 dimensions
CREATE TABLE game_embeddings (
    id bigserial PRIMARY KEY,
    game_id bigint,
    property_name text,
    embedding vector(384)
);

ALTER TABLE game_embeddings ADD CONSTRAINT fk_game_embedding FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE;
