## ADDED Requirements

### Requirement: Exercise Generation
The system SHALL generate various exercise types from learning content.

#### Scenario: Generate from article
- **WHEN** a user requests exercises from an article
- **THEN** the system SHALL use AI to generate:
  - Fill-in-the-blank questions (vocabulary)
  - Reading comprehension questions
  - Sentence reordering (optional)
- **AND** exercises SHALL be appropriate for user's level

#### Scenario: Exercise variety
- **WHEN** generating exercises
- **THEN** the system SHALL rotate exercise types
- **AND** avoid repetitive question formats

### Requirement: Fill-in-the-Blank Exercises
The system SHALL provide vocabulary-focused fill-in-the-blank exercises.

#### Scenario: Vocabulary cloze
- **WHEN** generating a fill-in-the-blank exercise
- **THEN** the system SHALL:
  - Select a key vocabulary word from the text
  - Remove it from the sentence
  - Provide context sufficient for inference
  - Optionally provide multiple choice options

#### Scenario: Answer validation
- **WHEN** a user submits an answer
- **THEN** the system SHALL:
  - Accept exact matches
  - Accept reasonable synonyms (with note)
  - Ignore case differences
  - Check spelling similarity for typos

### Requirement: Reading Comprehension Exercises
The system SHALL test understanding of article content.

#### Scenario: Comprehension questions
- **WHEN** generating comprehension exercises
- **THEN** questions SHALL:
  - Test main idea understanding
  - Test specific detail recall
  - Test inference ability
- **AND** answers SHALL be derivable from the text

#### Scenario: Question types
- **WHEN** displaying questions
- **THEN** formats MAY include:
  - Multiple choice (4 options)
  - True/False
  - Short answer (evaluated by AI)

#### Scenario: Passage reference
- **WHEN** displaying a comprehension question
- **THEN** the relevant passage SHALL be visible
- **OR** user can click to reveal the passage

### Requirement: Grammar Exercises
The system SHALL provide grammar-focused practice.

#### Scenario: Sentence correction
- **WHEN** generating grammar exercises
- **THEN** the system SHALL:
  - Present sentences with intentional errors
  - Ask user to identify and correct errors
  - Focus on common learner mistakes

#### Scenario: Grammar patterns
- **WHEN** targeting weak grammar areas
- **THEN** exercises SHALL focus on:
  - Patterns identified from user's error history
  - Common ESL trouble spots (articles, prepositions, tenses)

### Requirement: Exercise Session
The system SHALL provide structured exercise sessions.

#### Scenario: Start exercise session
- **WHEN** a user starts an exercise session
- **THEN** the system SHALL:
  - Present exercises one at a time
  - Track time per question
  - Allow skipping (with penalty note)
  - Show progress indicator

#### Scenario: Session completion
- **WHEN** an exercise session ends
- **THEN** the system SHALL display:
  - Score (percentage correct)
  - Time taken
  - Breakdown by question type
  - Detailed review of incorrect answers

#### Scenario: Exercise history
- **WHEN** viewing exercise history
- **THEN** users SHALL see:
  - Past exercise scores
  - Improvement trends
  - Questions frequently missed

### Requirement: Adaptive Difficulty
The system SHALL adjust exercise difficulty based on performance.

#### Scenario: Difficulty scaling
- **WHEN** generating exercises
- **THEN** difficulty SHALL match user's:
  - Assessed proficiency level
  - Recent exercise performance
  - Topic familiarity

#### Scenario: Challenge mode
- **WHEN** a user consistently scores >90%
- **THEN** the system MAY offer:
  - "Challenge" exercises at higher difficulty
  - Bonus vocabulary from advanced content
  - Timed challenge modes

#### Scenario: Support mode
- **WHEN** a user struggles (scores <60%)
- **THEN** the system SHALL:
  - Reduce difficulty temporarily
  - Provide more hints
  - Focus on foundational concepts
  - Offer encouragement

