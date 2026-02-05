## ADDED Requirements

### Requirement: Daily Learning Plan
The system SHALL generate a personalized daily learning plan for each user.

#### Scenario: Generate daily plan
- **WHEN** a user opens the dashboard
- **THEN** the system SHALL generate today's learning tasks based on:
  - User's proficiency level
  - User's daily time goal
  - Due vocabulary reviews (SRS)
  - Recommended new content
- **AND** display tasks in priority order

#### Scenario: Plan contents
- **WHEN** generating a daily plan
- **THEN** the plan SHALL include:
  - Review task (if vocabulary due for review)
  - Reading task (1 article at appropriate level)
  - Practice task (1-2 exercises targeting weak areas)
- **AND** each task SHALL show estimated time

#### Scenario: Task completion
- **WHEN** a user completes a task
- **THEN** the system SHALL mark it as done
- **AND** log the completion to progress tracking
- **AND** update the daily goal progress

### Requirement: Content Recommendation
The system SHALL recommend learning content tailored to user's level and interests.

#### Scenario: Filter by difficulty
- **WHEN** recommending articles
- **THEN** the system SHALL filter RSS content by estimated reading level
- **AND** prioritize content matching user's CEFR level (±1 level tolerance)

#### Scenario: Filter by interest
- **WHEN** recommending articles
- **THEN** the system SHALL prioritize content matching user's declared interests
- **AND** allow users to discover new topics occasionally (10% exploration)

#### Scenario: Difficulty tagging
- **WHEN** displaying article options
- **THEN** each article SHALL show:
  - Estimated difficulty (Easy/Medium/Hard or CEFR level)
  - Topic/category tag
  - Estimated reading time

### Requirement: Weekly Goals
The system SHALL set and track weekly learning goals.

#### Scenario: Auto-generate weekly goals
- **WHEN** a new week starts
- **THEN** the system SHALL generate goals based on:
  - User's daily time goal × 5 (assuming 5 active days)
  - Previous week's performance
  - Vocabulary retention targets

#### Scenario: Weekly goal examples
- **WHEN** displaying weekly goals
- **THEN** goals MAY include:
  - "Learn 20 new vocabulary words"
  - "Read 5 articles"
  - "Maintain 80% review accuracy"
  - "Complete 7-day streak"

#### Scenario: Goal adjustment
- **WHEN** a user consistently exceeds or misses goals
- **THEN** the system SHALL suggest goal adjustments
- **AND** allow users to manually adjust goals

### Requirement: Learning Path
The system SHALL provide a structured path for improvement.

#### Scenario: Identify weak areas
- **WHEN** analyzing user's progress
- **THEN** the system SHALL identify areas needing improvement:
  - Low vocabulary retention rate
  - Specific grammar patterns with high error rate
  - Reading comprehension gaps

#### Scenario: Targeted practice
- **WHEN** weak areas are identified
- **THEN** the daily plan SHALL include targeted exercises
- **AND** track improvement in those specific areas

