# Architecture Diagram

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                      English Learning Agent                      │
│                         (CLI Interface)                          │
└───────┬─────────────────────┬─────────────────────┬─────────────┘
        │                     │                     │
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌────────────────┐    ┌────────────────┐
│  RSS Feed     │    │    Content     │    │      AI        │
│   Reader      │    │   Extractor    │    │   Assistant    │
├───────────────┤    ├────────────────┤    ├────────────────┤
│ • Fetch feeds │    │ • Extract      │    │ • Simplify     │
│ • Parse HTML  │    │   phrases (x3) │    │   paragraphs   │
│ • Cache       │    │ • Identify     │    │ • Explain      │
│   articles    │    │   structures   │    │   articles     │
│ • Track read  │    │   (x1)         │    │ • Cache        │
│   status      │    │ • Rank by      │    │   responses    │
│               │    │   usefulness   │    │ • Track usage  │
└───────┬───────┘    └────────┬───────┘    └────────┬───────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
                     ┌────────▼────────┐
                     │  SQLite Storage │
                     ├─────────────────┤
                     │ • Articles      │
                     │ • Feeds         │
                     │ • Phrases       │
                     │ • Structures    │
                     │ • AI Cache      │
                     │ • Metadata      │
                     └─────────────────┘
```

## Daily Workflow Sequence

```
User                CLI              RSS Reader      Extractor       AI Assistant
 │                   │                    │              │                │
 │─────list────────>│                    │              │                │
 │                   │──fetch headlines──>│              │                │
 │                   │<──return 10-20─────│              │                │
 │<──show headlines──│                    │              │                │
 │                   │                    │              │                │
 │─────read 5──────>│                    │              │                │
 │                   │──get article───────>│              │                │
 │                   │<──article content───│              │                │
 │<──display article─│                    │              │                │
 │                   │                    │              │                │
 │────extract──────>│                    │              │                │
 │                   │──analyze text──────────────────>│                │
 │                   │<──3 phrases + 1 structure───────│                │
 │<──show results────│                    │              │                │
 │                   │                    │              │                │
 │─simplify para 2─>│                    │              │                │
 │                   │──simplify request──────────────────────────────>│
 │                   │<──simplified text────────────────────────────────│
 │<──show simplified─│                    │              │                │
 │                   │                    │              │                │
 │────explain──────>│                    │              │                │
 │                   │──explain request───────────────────────────────>│
 │                   │<──explanation──────────────────────────────────────│
 │<──show explain────│                    │              │                │
```

## Data Flow

```
External RSS Feeds
        │
        ▼
  ┌──────────┐
  │  Fetcher │ ──> Parse XML/Atom
  └────┬─────┘
       │
       ▼
  ┌──────────┐
  │  Parser  │ ──> Extract article URLs
  └────┬─────┘
       │
       ▼
  ┌──────────┐
  │ HTML     │ ──> Clean & extract main text
  │ Extractor│
  └────┬─────┘
       │
       ▼
  ┌──────────┐
  │ SQLite   │ <──> Store/Retrieve
  │ Database │
  └────┬─────┘
       │
       ├─────────────────────────┬──────────────────────┐
       ▼                         ▼                      ▼
  ┌──────────┐           ┌──────────┐          ┌──────────┐
  │  Phrase  │           │Sentence  │          │   AI     │
  │ Extractor│           │ Pattern  │          │ Client   │
  └────┬─────┘           └────┬─────┘          └────┬─────┘
       │                      │                     │
       └──────────────────────┴─────────────────────┘
                              │
                              ▼
                        User Display
```

## Component Dependencies

```
main.go
  │
  ├── cmd/
  │   ├── list.go     → rss.Reader
  │   ├── read.go     → rss.Reader + extractor.Analyzer
  │   ├── extract.go  → extractor.Analyzer
  │   ├── simplify.go → ai.Client
  │   ├── explain.go  → ai.Client
  │   └── review.go   → storage.DB
  │
  ├── internal/
  │   ├── rss/
  │   │   ├── reader.go    (fetch & parse)
  │   │   └── article.go   (HTML extraction)
  │   │
  │   ├── extractor/
  │   │   ├── phrases.go   (phrase extraction)
  │   │   └── structure.go (sentence patterns)
  │   │
  │   ├── ai/
  │   │   ├── client.go    (Eino integration)
  │   │   ├── prompts.go   (templates)
  │   │   └── cache.go     (response caching)
  │   │
  │   └── storage/
  │       ├── db.go        (SQLite wrapper)
  │       └── schema.go    (table definitions)
  │
  └── pkg/
      ├── config/
      │   └── config.go    (YAML/JSON loader)
      └── models/
          └── types.go     (shared data structures)
```

## Database Schema

```sql
-- Feeds table
CREATE TABLE feeds (
    id INTEGER PRIMARY KEY,
    url TEXT UNIQUE NOT NULL,
    title TEXT,
    last_fetched DATETIME,
    is_active BOOLEAN DEFAULT 1
);

-- Articles table
CREATE TABLE articles (
    id INTEGER PRIMARY KEY,
    feed_id INTEGER,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    content TEXT,
    published_at DATETIME,
    fetched_at DATETIME,
    is_read BOOLEAN DEFAULT 0,
    FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

-- Phrases table (extracted learning items)
CREATE TABLE phrases (
    id INTEGER PRIMARY KEY,
    article_id INTEGER,
    phrase TEXT NOT NULL,
    context TEXT,  -- sentence containing phrase
    extracted_at DATETIME,
    user_notes TEXT,
    FOREIGN KEY (article_id) REFERENCES articles(id)
);

-- Sentence structures table
CREATE TABLE structures (
    id INTEGER PRIMARY KEY,
    article_id INTEGER,
    sentence TEXT NOT NULL,
    pattern_type TEXT,  -- e.g., "relative clause", "conditional"
    extracted_at DATETIME,
    user_notes TEXT,
    FOREIGN KEY (article_id) REFERENCES articles(id)
);

-- AI cache table
CREATE TABLE ai_cache (
    id INTEGER PRIMARY KEY,
    request_hash TEXT UNIQUE NOT NULL,
    request_type TEXT,  -- "simplify" or "explain"
    response TEXT NOT NULL,
    created_at DATETIME,
    tokens_used INTEGER
);
```

