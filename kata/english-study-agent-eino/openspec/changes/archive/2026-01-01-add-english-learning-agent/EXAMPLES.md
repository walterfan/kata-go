# Usage Examples

This document provides concrete examples of how the English Learning Agent will work in practice.

---

## Initial Setup

```bash
# Install the agent
go install github.com/yourorg/english-learning-agent@latest

# Initialize with default configuration
english-agent init

# Output:
# ✓ Created config file: ~/.english-agent/config.yaml
# ✓ Initialized database: ~/.english-agent/data.db
# ✓ Added default feeds:
#   - BBC Learning English
#   - VOA Learning English
# 
# Configure your AI API key:
#   english-agent config set ai.api_key "your-key-here"
```

---

## Daily Workflow (5-10 minutes)

### Step 1: Skim Headlines

```bash
$ english-agent list

╔══════════════════════════════════════════════════════════════════╗
║               English Learning Agent - Headlines                  ║
╚══════════════════════════════════════════════════════════════════╝

[BBC Learning English] (5 unread)
  1. ⭐ The difference between 'affect' and 'effect'
     Published: 2 hours ago

  2. ⭐ 10 phrasal verbs for daily conversation
     Published: 1 day ago

  3. How to use present perfect tense
     Published: 2 days ago

[VOA Learning English] (3 unread)
  4. ⭐ News words: Technology terms explained
     Published: 3 hours ago

  5. American idioms: Weather expressions
     Published: 1 day ago

⭐ = Unread  |  Total: 8 articles  |  Time estimate: 5-10 min

Commands:
  • read <number>     - Open article for learning
  • refresh           - Fetch latest headlines
  • feeds             - Manage RSS feeds
```

### Step 2: Open One Article

```bash
$ english-agent read 1

╔══════════════════════════════════════════════════════════════════╗
║ The difference between 'affect' and 'effect'                     ║
║ Source: BBC Learning English | Published: 2 hours ago            ║
╚══════════════════════════════════════════════════════════════════╝

Many English learners confuse 'affect' and 'effect' because they sound 
similar. However, they have different meanings and uses.

'Affect' is usually a verb, meaning to influence or make a difference 
to something. For example: "The weather can affect your mood."

'Effect' is usually a noun, meaning the result of a change or influence.
For example: "The medicine had a positive effect on her health."

A good way to remember: Affect is the Action (both start with 'A'), 
and Effect is the End result (both start with 'E').

There are exceptions - 'effect' can be used as a verb in formal 
contexts, meaning to bring about or cause. For example: "The new 
policy will effect change in our organization."

Try this mnemonic: RAVEN - Remember, Affect is a Verb, Effect is a Noun.

[Article continues...]

Commands:
  • extract           - Extract phrases and structures
  • simplify <para>   - Simplify a paragraph
  • explain           - Get article explanation
  • next              - Next article
```

### Step 3: Extract Learning Items

```bash
$ english-agent extract

╔══════════════════════════════════════════════════════════════════╗
║ Extracted Learning Items - Article #1                            ║
╚══════════════════════════════════════════════════════════════════╝

📝 USEFUL PHRASES (3)

1. "make a difference to"
   Context: "...meaning to influence or make a difference to something."
   Type: Phrasal expression
   Usage: This phrase means to have an effect or impact on something.

2. "bring about"
   Context: "...'effect' can be used as a verb...meaning to bring about..."
   Type: Phrasal verb
   Usage: To cause something to happen; to make something occur.

3. "in formal contexts"
   Context: "...'effect' can be used as a verb in formal contexts..."
   Type: Prepositional phrase
   Usage: Used to indicate situations requiring professional language.

🔍 SENTENCE STRUCTURE (1)

Pattern: "A good way to remember: X is Y (both start with 'Z')"
Example: "A good way to remember: Affect is the Action (both start 
         with 'A'), and Effect is the End result (both start with 'E')."
Type: Memory aid pattern with parallel structure
Usage: This pattern is useful for teaching and explaining comparisons.

✓ Items saved to your learning collection

Commands:
  • review            - Review all saved items
  • export            - Export learning items
  • select            - Manually select phrases
```

### Step 4a: Simplify a Paragraph (Optional)

```bash
$ english-agent simplify "There are exceptions - 'effect' can be used as a verb in formal contexts, meaning to bring about or cause."

╔══════════════════════════════════════════════════════════════════╗
║ Text Simplification                                               ║
╚══════════════════════════════════════════════════════════════════╝

ORIGINAL:
"There are exceptions - 'effect' can be used as a verb in formal 
contexts, meaning to bring about or cause."

SIMPLIFIED (Intermediate Level):
"Sometimes 'effect' can be a verb. We use it in professional writing
to mean 'to cause something to happen'."

[Tokens used: 45 | Cached: No]

Commands:
  • simplify --level beginner  - Use simpler words
  • simplify --level advanced  - Keep more complexity
```

### Step 4b: Explain the Article (Optional)

```bash
$ english-agent explain

╔══════════════════════════════════════════════════════════════════╗
║ Article Explanation                                               ║
╚══════════════════════════════════════════════════════════════════╝

📚 SIMPLE EXPLANATION

This article teaches you the difference between two confusing English 
words: 'affect' and 'effect'.

KEY POINTS:
1. 'Affect' = verb (action word) = to influence something
   Example: Stress affects your health.

2. 'Effect' = noun (thing) = the result of something
   Example: The effect of stress is poor health.

3. Memory trick: Think of "RAVEN"
   - Affect = Verb (both have 'A' and 'V')
   - Effect = Noun (both have 'E' and 'N')

4. Advanced note: 'Effect' can sometimes be a verb in formal English,
   but you'll rarely need this.

💡 WHY THIS MATTERS FOR YOUR ENGLISH
Knowing this difference makes your writing clearer and more 
professional. Native speakers notice when these words are mixed up!

[Tokens used: 180 | Cached: No]
```

---

## Additional Commands

### Review Past Learning Items

```bash
$ english-agent review

╔══════════════════════════════════════════════════════════════════╗
║ Your Learning Collection                                          ║
╚══════════════════════════════════════════════════════════════════╝

📊 STATISTICS
Total phrases saved: 47
Total structures saved: 15
Articles read: 15
Learning streak: 5 days

📝 RECENT PHRASES (Last 7 days)
1. "make a difference to" (from: affect vs effect)
2. "bring about" (from: affect vs effect)
3. "in formal contexts" (from: affect vs effect)
4. "make sense of" (from: understanding idioms)
5. "come up with" (from: creative thinking)
...

🔍 RECENT STRUCTURES (Last 7 days)
1. Memory aid pattern (from: affect vs effect)
2. Conditional with 'if' clause (from: grammar tips)
3. Passive voice transformation (from: writing style)
...

Commands:
  • review --phrases      - Show only phrases
  • review --structures   - Show only structures
  • review --export       - Export to file
  • review --search <term> - Search your collection
```

### Manage Feeds

```bash
$ english-agent feeds

╔══════════════════════════════════════════════════════════════════╗
║ RSS Feed Management                                               ║
╚══════════════════════════════════════════════════════════════════╝

ACTIVE FEEDS (2)
1. BBC Learning English
   URL: https://www.bbc.co.uk/learningenglish/english/features/feed
   Last fetched: 2 hours ago
   Articles: 127 total, 5 unread

2. VOA Learning English
   URL: https://learningenglish.voanews.com/api/zvgove_kmeqq
   Last fetched: 2 hours ago
   Articles: 93 total, 3 unread

Commands:
  • feeds add <url>       - Add custom feed
  • feeds remove <id>     - Remove feed
  • feeds refresh         - Update all feeds
```

### Configure Settings

```bash
$ english-agent config

╔══════════════════════════════════════════════════════════════════╗
║ Configuration Settings                                            ║
╚══════════════════════════════════════════════════════════════════╝

AI SETTINGS
  Provider: openai
  Model: gpt-4
  API Key: sk-...Zx9w (configured ✓)
  Daily limit: 50 requests (12 used today)

EXTRACTION SETTINGS
  Phrases per article: 3
  Structures per article: 1
  Auto-extract: enabled

DISPLAY SETTINGS
  Headlines per page: 20
  Simplification level: intermediate

Commands:
  • config set <key> <value>  - Update setting
  • config get <key>          - View setting
  • config reset              - Reset to defaults
  • config show               - Show all settings
```

### Export Learning Items

```bash
$ english-agent review --export collection.json

Exporting your learning collection...
✓ Exported 47 phrases
✓ Exported 15 sentence structures
✓ Saved to: /Users/you/collection.json

You can also export as:
  • --export file.csv      (CSV format)
  • --export file.md       (Markdown format)
```

---

## Example Configuration File

**~/.english-agent/config.yaml**

```yaml
# AI Provider Configuration
ai:
  provider: openai
  api_key: ${OPENAI_API_KEY}  # or set directly
  model: gpt-4
  fallback_model: gpt-3.5-turbo
  timeout: 30s
  daily_request_limit: 50
  cache_responses: true
  cache_expiry_days: 30

# RSS Feeds
feeds:
  - url: https://www.bbc.co.uk/learningenglish/english/features/feed
    title: BBC Learning English
    enabled: true
  - url: https://learningenglish.voanews.com/api/zvgove_kmeqq
    title: VOA Learning English
    enabled: true
  # Add your custom feeds here
  # - url: https://example.com/feed.xml
  #   title: My Custom Feed
  #   enabled: true

# Extraction Settings
extraction:
  phrases_per_article: 3
  structures_per_article: 1
  auto_extract: true
  min_phrase_words: 2
  max_phrase_words: 5
  min_sentence_words: 10
  max_sentence_words: 25

# Display Settings
display:
  headlines_per_page: 20
  simplification_level: intermediate  # beginner | intermediate | advanced
  show_token_usage: true

# Storage
storage:
  database_path: ~/.english-agent/data.db
  auto_backup: true
  backup_frequency: weekly
```

---

## Error Handling Examples

### Network Error

```bash
$ english-agent list

⚠ Warning: Could not fetch feeds (network error)
  Showing cached headlines from 2 hours ago

[Cached headlines displayed...]
```

### AI Service Error

```bash
$ english-agent simplify "Some text..."

✗ Error: AI service unavailable
  The AI provider is not responding. This might be temporary.

Options:
  • Try again in a few minutes
  • Check your internet connection
  • Verify your API key: english-agent config get ai.api_key

Would you like to retry? [y/N]
```

### Invalid Configuration

```bash
$ english-agent config set ai.api_key "invalid"

Testing API key...
✗ Error: Invalid API key
  The provided key could not be authenticated.

Please get a valid API key from: https://platform.openai.com/api-keys
```

---

## Tips for Efficient Learning

### Daily Routine
```bash
# Morning routine (5 minutes)
english-agent list | head -10
english-agent read 1
english-agent extract

# Evening review (5 minutes)
english-agent review --phrases | tail -10
english-agent review --structures | tail -5
```

### Weekly Export
```bash
# Every Sunday, export your weekly collection
english-agent review --export weekly-$(date +%Y-%m-%d).json
```

### Batch Processing
```bash
# Read and extract from multiple articles
for i in {1..3}; do
  english-agent read $i --auto-extract
done
```

---

## Integration Examples

### Anki Flashcard Export (Future Enhancement)

```bash
$ english-agent review --export anki.csv --format anki

Exported 47 phrases in Anki-compatible format
Import into Anki: File → Import → anki.csv
```

### Daily Digest Email (Future Enhancement)

```bash
$ english-agent digest --email

Daily learning digest sent to: you@example.com
  • 3 new articles
  • 9 phrases extracted this week
  • 5-day learning streak!
```

---

This comprehensive example guide shows exactly how the English Learning Agent will work in practice, making it easy for users to understand the daily workflow and all available features.

