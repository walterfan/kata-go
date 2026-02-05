## ADDED Requirements

### Requirement: User Registration
The system SHALL allow users to create a personal learning profile.

#### Scenario: New user registration
- **WHEN** a new user accesses the system for the first time
- **THEN** the system SHALL prompt for basic profile information:
  - Username (required)
  - Native language (default: Chinese)
  - Daily learning goal in minutes (default: 10)
  - Areas of interest (technology, medicine, business, general)
- **AND** create a new user record with default settings

#### Scenario: Profile validation
- **WHEN** a user submits profile information
- **THEN** the system SHALL validate that username is unique
- **AND** reject duplicate usernames with a clear error message

### Requirement: User Profile Management
The system SHALL allow users to view and update their learning profile.

#### Scenario: View profile
- **WHEN** a user requests to view their profile
- **THEN** the system SHALL display:
  - Current proficiency level (if assessed)
  - Learning statistics (days active, vocabulary learned)
  - Current streak and longest streak
  - Areas of interest and preferences

#### Scenario: Update preferences
- **WHEN** a user updates their preferences
- **THEN** the system SHALL save:
  - Daily learning goal
  - Interest areas
  - Preferred difficulty level
- **AND** reflect changes in future content recommendations

### Requirement: Proficiency Assessment
The system SHALL assess user's English proficiency level through adaptive testing.

#### Scenario: Initial assessment
- **WHEN** a user completes profile setup
- **THEN** the system SHALL offer an optional proficiency assessment
- **AND** the assessment SHALL include 5-10 adaptive questions covering:
  - Vocabulary (word meaning in context)
  - Reading comprehension (short passage questions)
  - Grammar (sentence correction)

#### Scenario: Adaptive difficulty
- **WHEN** a user answers an assessment question
- **THEN** the system SHALL adjust the next question's difficulty:
  - Correct answer → increase difficulty
  - Incorrect answer → decrease difficulty
- **AND** converge on user's true level within 10 questions

#### Scenario: Assessment results
- **WHEN** a user completes the assessment
- **THEN** the system SHALL display:
  - Estimated CEFR level (A1-C2)
  - Breakdown by skill (vocabulary, reading, grammar)
  - Specific recommendations for improvement
- **AND** store the results in the user profile

### Requirement: Session Management
The system SHALL maintain user context across sessions.

#### Scenario: Remember user
- **WHEN** a user returns to the system
- **THEN** the system SHALL recognize them (via local storage or cookie)
- **AND** restore their learning context

#### Scenario: Switch user
- **WHEN** a user wants to switch accounts
- **THEN** the system SHALL allow logging out
- **AND** clear local session data

