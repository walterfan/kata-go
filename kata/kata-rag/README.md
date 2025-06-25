# Overview


```
[Source code files, docs]
      ↓
[Embedding model: BGE / Qwen]
      ↓ (text → vector)
[pgvector in Postgres]
      ↓ (similarity search)
[LLM (DeepSeek, GPT-4o)]
      ↓ (generate answer)
```


## setup

```
# If using docker image
docker run -d --name pgvector -p 5432:5432 \
  -e POSTGRES_PASSWORD=postgres \
  ankane/pgvector

# OR install in your local PostgreSQL:
CREATE EXTENSION vector;

# create table
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    path TEXT,
    content TEXT,
    embedding VECTOR(768)  -- if using nomic-embed-text embedding with 768 dims
);
```