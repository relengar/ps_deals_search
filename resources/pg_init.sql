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
    expiration timestamp,
    description_embedding vector(384)
);

