package storage

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dbPath string) {
	var err error
	if dbPath == "" {
		dbPath = "english_agent.db"
	}
	
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	createTables()
}

func createTables() {
	schema := `
	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		link TEXT UNIQUE,
		description TEXT,
		source TEXT,
		published_at TEXT,
		fetched_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS learning_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT, -- 'phrase' or 'structure'
		content TEXT,
		context TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS ai_cache (
		request_hash TEXT PRIMARY KEY,
		response TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS custom_feeds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		url TEXT UNIQUE NOT NULL,
		category TEXT DEFAULT '',
		enabled INTEGER DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := DB.Exec(schema)
	if err != nil {
		log.Fatalf("failed to create tables: %v", err)
	}
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

// Simple methods for saving/loading

func SaveArticle(title, link, desc, source, published string) error {
	_, err := DB.Exec(`INSERT OR IGNORE INTO articles (title, link, description, source, published_at) VALUES (?, ?, ?, ?, ?)`, 
		title, link, desc, source, published)
	return err
}

func SaveLearningItem(itemType, content, context string) error {
	_, err := DB.Exec(`INSERT INTO learning_items (type, content, context) VALUES (?, ?, ?)`, 
		itemType, content, context)
	return err
}

type LearningItem struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Context   string `json:"context"`
	CreatedAt string `json:"created_at"`
}

func GetLearningItems() ([]LearningItem, error) {
	rows, err := DB.Query(`SELECT id, type, content, context, created_at FROM learning_items ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []LearningItem
	for rows.Next() {
		var i LearningItem
		if err := rows.Scan(&i.ID, &i.Type, &i.Content, &i.Context, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

// Cache methods

func HashRequest(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func GetCache(hash string) (string, bool) {
	var response string
	err := DB.QueryRow(`SELECT response FROM ai_cache WHERE request_hash = ?`, hash).Scan(&response)
	if err != nil {
		return "", false
	}
	return response, true
}

func SetCache(hash, response string) error {
	_, err := DB.Exec(`INSERT OR REPLACE INTO ai_cache (request_hash, response) VALUES (?, ?)`, hash, response)
	return err
}

// Custom Feed CRUD operations

type CustomFeed struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Category  string `json:"category"`
	Enabled   bool   `json:"enabled"`
	CreatedAt string `json:"created_at"`
}

func AddCustomFeed(title, url, category string) error {
	_, err := DB.Exec(`INSERT INTO custom_feeds (title, url, category) VALUES (?, ?, ?)`,
		title, url, category)
	return err
}

func UpdateCustomFeed(id int, title, url, category string, enabled bool) error {
	enabledInt := 0
	if enabled {
		enabledInt = 1
	}
	_, err := DB.Exec(`UPDATE custom_feeds SET title = ?, url = ?, category = ?, enabled = ? WHERE id = ?`,
		title, url, category, enabledInt, id)
	return err
}

func DeleteCustomFeed(id int) error {
	_, err := DB.Exec(`DELETE FROM custom_feeds WHERE id = ?`, id)
	return err
}

func GetCustomFeeds() ([]CustomFeed, error) {
	rows, err := DB.Query(`SELECT id, title, url, category, enabled, created_at FROM custom_feeds ORDER BY category, title`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []CustomFeed
	for rows.Next() {
		var f CustomFeed
		var enabledInt int
		if err := rows.Scan(&f.ID, &f.Title, &f.URL, &f.Category, &enabledInt, &f.CreatedAt); err != nil {
			return nil, err
		}
		f.Enabled = enabledInt == 1
		feeds = append(feeds, f)
	}
	return feeds, nil
}

func GetEnabledCustomFeeds() ([]CustomFeed, error) {
	rows, err := DB.Query(`SELECT id, title, url, category, enabled, created_at FROM custom_feeds WHERE enabled = 1 ORDER BY category, title`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []CustomFeed
	for rows.Next() {
		var f CustomFeed
		var enabledInt int
		if err := rows.Scan(&f.ID, &f.Title, &f.URL, &f.Category, &enabledInt, &f.CreatedAt); err != nil {
			return nil, err
		}
		f.Enabled = enabledInt == 1
		feeds = append(feeds, f)
	}
	return feeds, nil
}
