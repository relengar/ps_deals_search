CREATE EXTENSION vector;

-- The embedding model we use (all-MiniLM-L6-v2) uses 384 dimensions
CREATE TABLE IF NOT EXISTS games (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    description text,
    price real NOT NULL,
    original_price real,
    url text NOT NULL,
    rating real NOT NULL,
    rating_sum int NOT NULL,
    expiration timestamp,
    platforms text[] default '{}'
);

-- The embedding model we use (all-MiniLM-L6-v2) uses 384 dimensions
CREATE TABLE IF NOT EXISTS game_embeddings (
    id bigserial PRIMARY KEY,
    game_id bigint NOT NULL,
    property_name text,
    embedding vector(384) NOT NULL
);

ALTER TABLE game_embeddings DROP CONSTRAINT IF EXISTS fk_game_embedding;
ALTER TABLE game_embeddings ADD CONSTRAINT fk_game_embedding FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE;
