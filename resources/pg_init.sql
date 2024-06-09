CREATE EXTENSION vector;

-- The embedding model we use (all-MiniLM-L6-v2) uses 384 dimensions
CREATE TABLE games (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    description text,
    price real NOT NULL,
    original_price real,
    url text NOT NULL,
    rating real NOT NULL,
    rating_sum int NOT NULL,
    expiration timestamp
);

-- The embedding model we use (all-MiniLM-L6-v2) uses 384 dimensions
CREATE TABLE game_embeddings (
    id bigserial PRIMARY KEY,
    game_id bigint NOT NULL,
    property_name text,
    embedding vector(384) NOT NULL
);

ALTER TABLE game_embeddings ADD CONSTRAINT fk_game_embedding FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE;
