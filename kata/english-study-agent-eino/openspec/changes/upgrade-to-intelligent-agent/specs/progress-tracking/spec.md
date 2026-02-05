## ADDED Requirements

### Requirement: Activity Logging
The system SHALL log all learning activities for progress analysis.

#### Scenario: Log reading activity
- **WHEN** a user reads an article
- **THEN** the system SHALL log:
  - Article identifier
  - Time spent reading
  - Actions taken (explain, translate, extract vocabulary)
  - Timestamp

#### Scenario: Log exercise activity
- **WHEN** a user completes an exercise
- **THEN** the system SHALL log:
  - Exercise type
  - Score/accuracy
  - Time spent
  - Specific questions and answers

#### Scenario: Log review activity
- **WHEN** a user completes vocabulary reviews
- **THEN** the system SHALL log:
  - Words reviewed
  - Recall accuracy
  - Session duration

### Requirement: Progress Dashboard
The system SHALL provide a visual dashboard of learning progress.

#### Scenario: Overview metrics
- **WHEN** a user views the dashboard
- **THEN** the system SHALL display key metrics:
  - Total vocabulary learned
  - Current streak (consecutive days)
  - This week's learning time
  - Overall progress percentage toward goals

#### Scenario: Vocabulary growth chart
- **WHEN** viewing progress
- **THEN** the system SHALL show a line chart of:
  - Words learned over time (cumulative)
  - Words reviewed and retained
  - New words vs. reviewed words per day

#### Scenario: Time investment chart
- **WHEN** viewing progress
- **THEN** the system SHALL show:
  - Daily learning time (bar chart, last 7/30 days)
  - Comparison to goal (highlighted goal line)

### Requirement: Streak Tracking
The system SHALL track learning streaks to encourage consistency.

#### Scenario: Update streak
- **WHEN** a user completes any learning activity
- **THEN** the system SHALL:
  - Mark today as an active day
  - Increment current streak if consecutive with yesterday
  - Reset streak to 1 if gap detected

#### Scenario: Display streak
- **WHEN** displaying streak
- **THEN** the system SHALL show:
  - Current streak count with 🔥 icon
  - Longest streak achieved
  - Calendar heatmap of activity (last 30/90 days)

#### Scenario: Streak milestones
- **WHEN** a user reaches streak milestones (7, 30, 100 days)
- **THEN** the system SHALL display a celebration notification
- **AND** award a streak badge

### Requirement: Progress Reports
The system SHALL generate periodic progress summaries.

#### Scenario: Weekly report
- **WHEN** a week ends
- **THEN** the system SHALL generate a summary:
  - Words learned this week
  - Time spent learning
  - Streak status
  - Comparison to previous week
  - Top achievements

#### Scenario: Monthly review
- **WHEN** viewing monthly progress
- **THEN** the system SHALL show:
  - Vocabulary growth (start of month vs now)
  - Goal completion rate
  - Strongest and weakest areas
  - Recommendations for next month

### Requirement: Comparative Analytics
The system SHALL provide insights by comparing current vs past performance.

#### Scenario: Show improvement
- **WHEN** a user views progress
- **THEN** the system SHALL highlight improvements:
  - "Your vocabulary retention improved by 15% this month"
  - "You're learning 20% faster than last week"

#### Scenario: Identify plateaus
- **WHEN** progress stagnates
- **THEN** the system SHALL:
  - Detect lack of improvement
  - Suggest strategy changes
  - Recommend varying content or exercises

