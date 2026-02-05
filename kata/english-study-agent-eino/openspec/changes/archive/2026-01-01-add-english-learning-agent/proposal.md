# Change: Add English Learning AI Agent with RSS Feed Integration

## Why
English learners need a time-efficient way to improve their vocabulary and comprehension by reading authentic content. Reading RSS feeds daily (5-10 minutes) with AI assistance helps learners extract useful phrases, understand sentence structures, and get simplified explanations, making language learning practical and sustainable.

## What Changes
- Add RSS feed reader to fetch and display English learning content headlines
- Add content extraction capability to identify useful phrases and sentence structures
- Add AI assistant integration to rewrite paragraphs in simpler English and explain content
- Add daily workflow that supports:
  - Skimming headlines (5-10 min/day)
  - Opening one article for deep reading
  - Extracting 3 useful phrases and 1 sentence structure
  - Optional AI-powered simplification and explanation

## Impact
- Affected specs: 
  - `rss-feed-reader` (new capability)
  - `content-extractor` (new capability)
  - `ai-assistant` (new capability)
- Affected code: 
  - New Go project structure with Eino framework
  - RSS parsing and content fetching modules
  - NLP-based phrase and sentence extraction
  - AI integration for text simplification

