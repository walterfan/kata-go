package rss

import (
	"context"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/walterfan/english-agent/internal/config"
	"golang.org/x/net/html"
)

type Article struct {
	Title       string
	Link        string
	Description string
	Published   string
	Source      string
}

type Fetcher struct {
	parser *gofeed.Parser
}

func NewFetcher() *Fetcher {
	fp := gofeed.NewParser()
	fp.Client = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &userAgentTransport{
			Transport: http.DefaultTransport,
		},
	}
	return &Fetcher{
		parser: fp,
	}
}

type userAgentTransport struct {
	Transport http.RoundTripper
}

func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	return t.Transport.RoundTrip(req)
}

func (f *Fetcher) FetchHeadlines(ctx context.Context) ([]Article, error) {
	return f.FetchFromSource(ctx, "")
}

// FetchFromSource fetches articles from a specific source or all sources if sourceTitle is empty
func (f *Fetcher) FetchFromSource(ctx context.Context, sourceTitle string) ([]Article, error) {
	cfg := config.Get()
	articles := []Article{} // Initialize as empty slice to return [] instead of null

	for _, feedCfg := range cfg.Feeds {
		// Skip if source filter is set and doesn't match
		if sourceTitle != "" && sourceTitle != "all" && feedCfg.Title != sourceTitle {
			continue
		}

		// Use a timeout context for each feed
		feedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		feed, err := f.parser.ParseURLWithContext(feedCfg.URL, feedCtx)
		if err != nil {
			// Log error but continue with other feeds
			continue
		}

		for _, item := range feed.Items {
			articles = append(articles, Article{
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				Published:   item.Published,
				Source:      feedCfg.Title,
			})
		}
	}

	return articles, nil
}

// GetSources returns the list of configured RSS sources
func (f *Fetcher) GetSources() []config.FeedConfig {
	return config.Get().Feeds
}

// WebArticle represents an article fetched from a web URL
type WebArticle struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	URL     string `json:"url"`
}

// FetchFromURL fetches article content from a web URL
func (f *Fetcher) FetchFromURL(ctx context.Context, url string) (*WebArticle, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set User-Agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	content := string(body)
	contentType := resp.Header.Get("Content-Type")
	
	var title, articleContent string
	
	// Check if it's HTML or plain text/markdown
	if strings.Contains(contentType, "text/html") || strings.HasPrefix(content, "<!") || strings.HasPrefix(content, "<html") {
		// HTML content - parse it
		title = extractTitle(content)
		articleContent = extractMainContent(content)
	} else {
		// Plain text or Markdown - use directly
		title = extractTitleFromText(content, url)
		articleContent = cleanMarkdown(content)
	}

	return &WebArticle{
		Title:   title,
		Content: articleContent,
		URL:     url,
	}, nil
}

// extractTitleFromText extracts title from plain text or markdown
func extractTitleFromText(content, url string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Check for markdown heading
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
		// First non-empty line as fallback
		if len(line) > 0 && len(line) < 200 {
			return line
		}
	}
	// Use URL as last resort
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "Untitled"
}

// cleanMarkdown removes markdown formatting for cleaner reading
func cleanMarkdown(content string) string {
	// Remove common markdown elements that aren't useful for reading
	lines := strings.Split(content, "\n")
	var result []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Skip empty lines at the start
		if len(result) == 0 && trimmed == "" {
			continue
		}
		
		// Skip image references
		if strings.HasPrefix(trimmed, "![") || strings.HasPrefix(trimmed, "<img") {
			continue
		}
		
		// Skip badge/shield references
		if strings.Contains(trimmed, "shields.io") || strings.Contains(trimmed, "badge") {
			continue
		}
		
		// Clean up headers (remove #)
		if strings.HasPrefix(trimmed, "#") {
			trimmed = strings.TrimLeft(trimmed, "# ")
			trimmed = "\n" + trimmed + "\n" // Add spacing around headers
		}
		
		// Clean up bold/italic
		trimmed = strings.ReplaceAll(trimmed, "**", "")
		trimmed = strings.ReplaceAll(trimmed, "__", "")
		
		// Clean up code blocks markers
		if strings.HasPrefix(trimmed, "```") {
			continue
		}
		
		// Clean up inline code
		re := regexp.MustCompile("`([^`]+)`")
		trimmed = re.ReplaceAllString(trimmed, "$1")
		
		// Clean up links [text](url) -> text
		linkRe := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
		trimmed = linkRe.ReplaceAllString(trimmed, "$1")
		
		result = append(result, trimmed)
	}
	
	return strings.TrimSpace(strings.Join(result, "\n"))
}

// extractTitle extracts the title from HTML
func extractTitle(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "Untitled"
	}

	var title string
	var findTitle func(*html.Node)
	findTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTitle(c)
		}
	}
	findTitle(doc)

	if title == "" {
		// Try to find h1
		findTitle = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "h1" {
				title = extractText(n)
			}
			if title == "" {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					findTitle(c)
				}
			}
		}
		findTitle(doc)
	}

	return strings.TrimSpace(title)
}

// extractMainContent extracts readable text from HTML
func extractMainContent(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return cleanText(htmlContent)
	}

	// Try to find article or main content
	var content string
	var findContent func(*html.Node) bool
	findContent = func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			// Look for article, main, or content divs
			if n.Data == "article" || n.Data == "main" {
				content = extractText(n)
				return true
			}
			// Check for common content class names
			for _, attr := range n.Attr {
				if attr.Key == "class" || attr.Key == "id" {
					if strings.Contains(attr.Val, "content") ||
						strings.Contains(attr.Val, "article") ||
						strings.Contains(attr.Val, "post-body") ||
						strings.Contains(attr.Val, "entry-content") {
						content = extractText(n)
						return true
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if findContent(c) {
				return true
			}
		}
		return false
	}
	findContent(doc)

	// Fallback: extract all paragraph text
	if content == "" {
		var paragraphs []string
		var findParagraphs func(*html.Node)
		findParagraphs = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "p" {
				text := extractText(n)
				if len(text) > 50 { // Only meaningful paragraphs
					paragraphs = append(paragraphs, text)
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				findParagraphs(c)
			}
		}
		findParagraphs(doc)
		content = strings.Join(paragraphs, "\n\n")
	}

	return cleanText(content)
}

// extractText extracts all text from a node
func extractText(n *html.Node) string {
	var text strings.Builder
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		// Skip script, style, nav, header, footer
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "nav", "header", "footer", "aside", "noscript":
				return
			}
		}
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
			text.WriteString(" ")
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(n)
	return text.String()
}

// cleanText cleans up extracted text
func cleanText(text string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	// Remove common noise
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.ReplaceAll(text, "\t", " ")

	// Trim
	text = strings.TrimSpace(text)

	return text
}

