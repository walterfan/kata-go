# Content Extractor Capability

## Requirements

### Requirement: Phrase Extraction
The system SHALL automatically extract useful English phrases from article content to support vocabulary learning.

#### Scenario: Extract common phrases
- **WHEN** a user opens an article for learning
- **THEN** the system SHALL identify and extract at least 3 useful phrases
- **AND** prioritize phrases that are:
  - Common in English (not overly technical or rare)
  - 2-5 words in length
  - Grammatically complete expressions (phrasal verbs, collocations, idioms)

#### Scenario: Rank by usefulness
- **WHEN** multiple phrases are identified
- **THEN** the system SHALL rank them by usefulness based on:
  - Frequency in general English corpus
  - Phrase completeness
  - Educational value (collocations, idioms)

#### Scenario: Display phrase context
- **WHEN** extracted phrases are displayed
- **THEN** each phrase SHALL be shown with:
  - The original sentence containing the phrase
  - The phrase highlighted within the sentence
  - A suggested definition or usage note (optional, if available)

### Requirement: Sentence Structure Identification
The system SHALL identify interesting sentence structures from articles to help learners understand English grammar patterns.

#### Scenario: Extract notable structures
- **WHEN** a user opens an article for learning
- **THEN** the system SHALL identify at least 1 sentence with an interesting grammatical structure
- **AND** prioritize sentences that demonstrate:
  - Common but sophisticated patterns (relative clauses, conditionals, passive voice)
  - Clear, well-formed structure
  - Moderate length (10-25 words)

#### Scenario: Highlight structure pattern
- **WHEN** a sentence structure is displayed
- **THEN** the system SHALL show:
  - The full sentence
  - A brief description of the grammatical pattern (e.g., "Relative clause with 'which'")
  - The structural components highlighted or labeled (optional)

#### Scenario: User manual selection
- **WHEN** a user disagrees with automatic extraction
- **THEN** the system SHALL allow manual selection of any sentence in the article
- **AND** allow marking it as a "favorite structure" for review

### Requirement: Learning Item Storage
The system SHALL store extracted phrases and sentence structures for later review and practice.

#### Scenario: Save extracted items
- **WHEN** phrases and structures are extracted
- **THEN** the system SHALL save them to local storage with:
  - The phrase or sentence text
  - The source article title and URL
  - The extraction date
  - User notes (optional)

#### Scenario: Review past extractions
- **WHEN** a user requests to review learning items
- **THEN** the system SHALL display all saved phrases and structures
- **AND** support filtering by date, article, or keyword

#### Scenario: Export learning items
- **WHEN** a user wants to export their collection
- **THEN** the system SHALL provide an export in a standard format (JSON, CSV, or markdown)

### Requirement: Content Analysis Quality
The system SHALL provide reasonably accurate content extraction with graceful handling of edge cases.

#### Scenario: Handle short articles
- **WHEN** an article is too short (< 100 words)
- **THEN** the system SHALL extract what it can
- **AND** notify the user that results may be limited

#### Scenario: Handle non-English content
- **WHEN** article content is not primarily in English
- **THEN** the system SHALL skip extraction
- **AND** warn the user that the article language may not match

#### Scenario: Handle extraction failure
- **WHEN** automatic extraction fails due to parsing errors
- **THEN** the system SHALL allow manual phrase and sentence selection
- **AND** log the error for debugging

