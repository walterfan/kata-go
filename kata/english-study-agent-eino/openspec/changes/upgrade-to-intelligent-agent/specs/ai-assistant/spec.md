## MODIFIED Requirements

### Requirement: Content Explanation
The system SHALL provide AI-generated explanations of article content tailored to English learners, **personalized to their proficiency level**.

#### Scenario: Explain article like to a learner
- **WHEN** a user requests an article explanation
- **THEN** the system SHALL send the article to an AI model with a prompt like "Explain this article like I'm a junior developer" (or similar learner-appropriate persona)
- **AND** display a concise summary with key points
- **AND** adjust complexity based on user's assessed CEFR level

#### Scenario: Explain specific phrases
- **WHEN** a user selects an unfamiliar phrase
- **THEN** the system SHALL request an explanation from the AI
- **AND** display:
  - The phrase meaning in simple terms
  - Example usage in a sentence
  - Similar phrases or synonyms (optional)
- **AND** optionally add the phrase to user's vocabulary queue

#### Scenario: Contextual explanation
- **WHEN** explaining content
- **THEN** the AI SHALL consider:
  - The learner's **actual** proficiency level (from assessment)
  - The article's subject domain
  - Cultural or idiomatic context where relevant
  - User's native language for comparison examples

#### Scenario: Personalized vocabulary extraction
- **WHEN** extracting vocabulary from an article
- **THEN** the system SHALL:
  - Filter words based on user's known vocabulary (avoid redundant suggestions)
  - Prioritize words slightly above user's current level
  - Automatically offer to add extracted words to SRS queue

## ADDED Requirements

### Requirement: Assessment Chain Integration
The system SHALL use AI to generate and evaluate proficiency assessment questions.

#### Scenario: Generate assessment question
- **WHEN** starting or continuing an assessment
- **THEN** the system SHALL use AI to generate:
  - A question appropriate for the current estimated level
  - Clear answer options (for multiple choice)
  - An explanation for the correct answer (stored for feedback)

#### Scenario: Evaluate assessment answer
- **WHEN** a user submits an assessment answer
- **THEN** the system SHALL:
  - Determine correctness
  - Adjust estimated level based on response
  - Select next question difficulty accordingly

### Requirement: Learning Plan Generation
The system SHALL use AI to generate personalized learning plans.

#### Scenario: Daily plan generation
- **WHEN** generating a daily learning plan
- **THEN** the system SHALL use AI to:
  - Analyze user's profile and recent progress
  - Select appropriate content from RSS feeds
  - Create a balanced plan (review + new content + practice)
  - Output structured task list

#### Scenario: Adaptive recommendations
- **WHEN** recommending content
- **THEN** the AI SHALL consider:
  - User's declared interests
  - Historical reading patterns
  - Identified weak areas
  - Variety (avoid repetitive topics)

### Requirement: Exercise Generation Chain
The system SHALL use AI to generate exercises from learning content.

#### Scenario: Generate fill-in-blank
- **WHEN** creating vocabulary exercises
- **THEN** the AI SHALL:
  - Select contextually rich sentences
  - Remove words appropriate to user's level
  - Generate distractor options (for multiple choice)

#### Scenario: Generate comprehension questions
- **WHEN** creating reading exercises
- **THEN** the AI SHALL:
  - Create questions testing different comprehension levels
  - Include main idea, detail, and inference questions
  - Ensure answers are derivable from the text

### Requirement: Intelligent Feedback Chain
The system SHALL use AI to provide personalized feedback on user responses.

#### Scenario: Exercise feedback
- **WHEN** evaluating an exercise answer
- **THEN** the AI SHALL generate:
  - Clear indication of correctness
  - Explanation of why the correct answer is correct
  - Encouragement appropriate to the situation
  - Suggestions for improvement if needed

#### Scenario: Progress feedback
- **WHEN** generating progress summaries
- **THEN** the AI SHALL:
  - Highlight achievements and improvements
  - Identify areas needing more attention
  - Provide actionable recommendations
  - Maintain an encouraging, supportive tone

