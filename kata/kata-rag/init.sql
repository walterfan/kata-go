CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    path TEXT,
    content TEXT,
    embedding VECTOR(768)  -- adjust dimension to match your embedding
);

