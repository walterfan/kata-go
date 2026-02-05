## ADDED Requirements

### Requirement: Text Simplification
The system SHALL provide AI-powered text simplification to help learners understand complex paragraphs.

#### Scenario: Simplify selected paragraph
- **WHEN** a user selects a paragraph and requests simplification
- **THEN** the system SHALL send the text to an AI model with a prompt like "Rewrite this paragraph in simpler English"
- **AND** display the simplified version alongside the original

#### Scenario: Adjust simplification level
- **WHEN** a user requests a specific simplification level
- **THEN** the system SHALL support levels like:
  - "Beginner" (elementary vocabulary, simple sentences)
  - "Intermediate" (common vocabulary, moderate complexity)
  - "Advanced" (academic vocabulary, preserved nuance)

#### Scenario: Handle simplification errors
- **WHEN** the AI service is unavailable or returns an error
- **THEN** the system SHALL display a user-friendly error message
- **AND** offer to retry or skip the operation

### Requirement: Content Explanation
The system SHALL provide AI-generated explanations of article content tailored to English learners.

#### Scenario: Explain article like to a learner
- **WHEN** a user requests an article explanation
- **THEN** the system SHALL send the article to an AI model with a prompt like "Explain this article like I'm a junior developer" (or similar learner-appropriate persona)
- **AND** display a concise summary with key points

#### Scenario: Explain specific phrases
- **WHEN** a user selects an unfamiliar phrase
- **THEN** the system SHALL request an explanation from the AI
- **AND** display:
  - The phrase meaning in simple terms
  - Example usage in a sentence
  - Similar phrases or synonyms (optional)

#### Scenario: Contextual explanation
- **WHEN** explaining content
- **THEN** the AI SHALL consider:
  - The learner's assumed proficiency level (intermediate by default)
  - The article's subject domain
  - Cultural or idiomatic context where relevant

### Requirement: AI Configuration
The system SHALL allow configuration of AI service parameters for flexibility and cost control.

#### Scenario: Configure AI provider
- **WHEN** the system initializes
- **THEN** it SHALL support configuration of:
  - AI provider (OpenAI, compatible APIs)
  - API endpoint and authentication
  - Model selection (e.g., gpt-4, gpt-3.5-turbo)

#### Scenario: Set request limits
- **WHEN** a user configures AI settings
- **THEN** the system SHALL allow setting:
  - Maximum tokens per request
  - Request timeout duration
  - Daily request limit (cost control)

#### Scenario: API key validation
- **WHEN** a user provides an API key
- **THEN** the system SHALL validate it with a test request
- **AND** display confirmation or error message

### Requirement: AI Response Caching
The system SHALL cache AI responses to reduce costs and improve performance for repeated queries.

#### Scenario: Cache simplification results
- **WHEN** a paragraph is simplified
- **THEN** the system SHALL cache the result keyed by original text and simplification level
- **AND** reuse the cached result if the same request is made again

#### Scenario: Cache explanation results
- **WHEN** content or phrases are explained
- **THEN** the system SHALL cache explanations
- **AND** display cached results instantly on subsequent requests

#### Scenario: Cache expiration
- **WHEN** the cache grows large or entries age
- **THEN** the system SHALL expire entries older than 30 days
- **AND** limit cache size to a configurable maximum (e.g., 1000 entries)

### Requirement: User Control and Transparency
The system SHALL give users control over AI usage and make AI operations transparent.

#### Scenario: Optional AI usage
- **WHEN** using the learning workflow
- **THEN** AI-powered features (simplification, explanation) SHALL be clearly marked as optional
- **AND** the core workflow SHALL function without AI if not configured

#### Scenario: Show AI processing indicator
- **WHEN** an AI request is in progress
- **THEN** the system SHALL display a loading indicator or progress message
- **AND** show estimated time if available

#### Scenario: Display token usage
- **WHEN** an AI operation completes
- **THEN** the system MAY display approximate token usage and cost (if configured)
- **AND** track cumulative daily usage against limits

