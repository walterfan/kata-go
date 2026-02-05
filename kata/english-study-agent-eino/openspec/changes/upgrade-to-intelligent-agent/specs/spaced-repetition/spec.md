## ADDED Requirements

### Requirement: Vocabulary Queue
The system SHALL maintain a vocabulary review queue using spaced repetition.

#### Scenario: Add vocabulary to queue
- **WHEN** a user extracts vocabulary from an article
- **THEN** the system SHALL:
  - Add each word/phrase to the review queue
  - Set initial review date to tomorrow
  - Store context (original sentence) and translation

#### Scenario: Manual vocabulary add
- **WHEN** a user manually adds a word
- **THEN** the system SHALL:
  - Create a review entry with user-provided word and meaning
  - Optionally allow adding context sentence
  - Schedule first review for next day

### Requirement: SM-2 Spaced Repetition
The system SHALL use a simplified SM-2 algorithm for scheduling reviews.

#### Scenario: Calculate next review
- **WHEN** a user reviews a vocabulary item
- **THEN** the system SHALL calculate next review based on recall quality:
  - Perfect recall (5): interval × ease_factor, increase ease
  - Good recall (4): interval × ease_factor
  - Okay recall (3): interval × ease_factor, decrease ease slightly
  - Difficult recall (2): reset interval to 1 day
  - Failed recall (0-1): review again in same session

#### Scenario: Initial scheduling
- **WHEN** a word is first added
- **THEN** default parameters SHALL be:
  - Ease factor: 2.5
  - First interval: 1 day
  - Second interval: 6 days (if first review successful)

#### Scenario: Ease factor bounds
- **WHEN** updating ease factor
- **THEN** ease factor SHALL:
  - Never go below 1.3
  - Increase by 0.15 for perfect recall
  - Decrease by 0.2 for difficult recall

### Requirement: Review Session
The system SHALL provide a dedicated review interface.

#### Scenario: Start review session
- **WHEN** a user has vocabulary due for review
- **THEN** the system SHALL:
  - Display count of due items
  - Show "Review Now" button
  - Estimate session duration (2-3 seconds per item)

#### Scenario: Review card display
- **WHEN** reviewing a vocabulary item
- **THEN** the system SHALL display:
  - The word/phrase
  - Context sentence (with target word hidden initially)
  - Option to reveal meaning
  - Self-assessment buttons (Again/Hard/Good/Easy)

#### Scenario: Review completion
- **WHEN** a review session ends
- **THEN** the system SHALL display:
  - Items reviewed count
  - Accuracy rate
  - Next review preview (items due tomorrow, this week)

### Requirement: Review Statistics
The system SHALL track vocabulary retention metrics.

#### Scenario: Retention rate
- **WHEN** viewing vocabulary statistics
- **THEN** the system SHALL show:
  - Overall retention rate (items recalled successfully / total reviews)
  - Retention by time period (7-day, 30-day, all-time)
  - Vocabulary mastery distribution (learning, young, mature)

#### Scenario: Vocabulary categories
- **WHEN** managing vocabulary
- **THEN** items SHALL be categorized:
  - Learning (0-2 successful reviews)
  - Young (3-6 successful reviews)
  - Mature (7+ successful reviews)

#### Scenario: Difficult words
- **WHEN** a word has low ease factor (<2.0) or many lapses
- **THEN** the system SHALL:
  - Flag it as a "difficult word"
  - Suggest additional study strategies
  - Prioritize it in exercises

### Requirement: Vocabulary Management
The system SHALL allow users to manage their vocabulary queue.

#### Scenario: View vocabulary list
- **WHEN** a user views their vocabulary
- **THEN** the system SHALL display:
  - Total words in queue
  - Words by status (due today, overdue, upcoming)
  - Search and filter options

#### Scenario: Edit vocabulary
- **WHEN** a user edits a vocabulary entry
- **THEN** they SHALL be able to modify:
  - Word/phrase
  - Meaning/translation
  - Context sentence
  - Personal notes

#### Scenario: Remove vocabulary
- **WHEN** a user removes a word
- **THEN** the system SHALL:
  - Delete it from the review queue
  - Retain it in history for analytics
  - Confirm before permanent deletion

