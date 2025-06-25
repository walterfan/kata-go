package vectorstore

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    "rag-assistant/internal/embedder"
)

func InsertEmbedding(ctx context.Context, db *sql.DB, path string, content string, embedding []float32) error {
    _, err := db.ExecContext(ctx, `
        INSERT INTO documents (path, content, embedding)
        VALUES ($1, $2, $3)
    `, path, content, embedder.VectorToPGArray(embedding))
    return err
}

func SearchSimilar(ctx context.Context, db *sql.DB, embedding []float32, topK int) ([]string, error) {
    rows, err := db.QueryContext(ctx, `
        SELECT path, content, embedding <-> $1 AS distance
        FROM documents
        ORDER BY embedding <-> $1
        LIMIT $2
    `, embedder.VectorToPGArray(embedding), topK)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []string
    for rows.Next() {
        var path string
        var content string
        var distance float64
        if err := rows.Scan(&path, &content, &distance); err != nil {
            return nil, err
        }
        const maxLines = 1000

        lines := strings.Split(content, "\n")
        var snippet string
        if len(lines) <= maxLines {
            snippet = strings.Join(lines, "\n")
        } else {
            snippet = strings.Join(lines[:maxLines], "\n") + "\n..."
        }

        results = append(results, fmt.Sprintf("📄 %s\n%s", path, snippet))
    }
    return results, nil
}
