# Tasks: Upgrade to Intelligent English Learning Agent

> **Strategy**: Implement phase by phase. Each phase must pass all acceptance tests before proceeding to the next.

---

## Phase 1: User Profile & Assessment (Week 1)

### 1.1 Database Foundation
- [ ] 1.1.1 Create migration script for new tables (`users`, `user_profiles`)
- [ ] 1.1.2 Add user-related DB functions to `internal/storage/user.go`
- [ ] 1.1.3 Add default user migration for existing data

### 1.2 User Management API
- [ ] 1.2.1 Add `POST /api/users/register` endpoint
- [ ] 1.2.2 Add `GET /api/users/profile` endpoint  
- [ ] 1.2.3 Add `PUT /api/users/profile` endpoint
- [ ] 1.2.4 Add simple session management (cookie-based)

### 1.3 Assessment System
- [ ] 1.3.1 Create assessment prompt templates in `internal/agent/assessment.go`
- [ ] 1.3.2 Add `POST /api/assessment/start` - generate first question
- [ ] 1.3.3 Add `POST /api/assessment/answer` - evaluate and get next question
- [ ] 1.3.4 Add `GET /api/assessment/result` - summarize assessment

### 1.4 UI: Profile & Assessment
- [ ] 1.4.1 Add welcome/onboarding page for new users
- [ ] 1.4.2 Add profile setup wizard (3 questions)
- [ ] 1.4.3 Add assessment quiz interface (5-10 questions)
- [ ] 1.4.4 Display assessment results with CEFR level
- [ ] 1.4.5 Add profile view/edit in settings

### ✅ Phase 1 Acceptance Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| P1-T1 | **New User Registration** | User can enter username, select native language, set daily goal |
| P1-T2 | **Profile Persistence** | After restart, user profile is retained |
| P1-T3 | **Assessment Start** | System generates vocabulary question appropriate for default level |
| P1-T4 | **Adaptive Assessment** | After correct answer, next question is harder; after wrong, easier |
| P1-T5 | **Assessment Completion** | After 5-10 questions, system shows CEFR level (A1-C2) |
| P1-T6 | **Profile Update** | User can change daily goal and interests, changes persist |
| P1-T7 | **Session Memory** | Returning user is recognized without re-registration |

**Phase 1 Definition of Done:**
- [ ] All P1-T1 to P1-T7 tests pass
- [ ] No regression in existing features (RSS, AI explain, etc.)

---

## Phase 2: Learning Plan Generation (Week 1-2)

### 2.1 Planning Backend
- [ ] 2.1.1 Create `learning_goals` table
- [ ] 2.1.2 Implement planner chain in `internal/agent/planner.go`
- [ ] 2.1.3 Add `GET /api/plan/today` - generate daily tasks
- [ ] 2.1.4 Add `POST /api/plan/complete-task` - mark task done

### 2.2 Content Recommendation
- [ ] 2.2.1 Add difficulty estimation for articles
- [ ] 2.2.2 Filter articles by user level and interests
- [ ] 2.2.3 Add `GET /api/articles/recommended` endpoint

### 2.3 UI: Daily Dashboard
- [ ] 2.3.1 Create dashboard showing today's tasks
- [ ] 2.3.2 Add task completion checkboxes
- [ ] 2.3.3 Show recommended articles with difficulty badges
- [ ] 2.3.4 Add daily progress indicator

### ✅ Phase 2 Acceptance Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| P2-T1 | **Daily Plan Generation** | On dashboard load, system shows 2-3 daily tasks |
| P2-T2 | **Plan Based on Level** | Beginner sees easier content, Advanced sees harder |
| P2-T3 | **Plan Based on Interests** | Tech-interested user gets tech articles first |
| P2-T4 | **Task Completion** | Checking a task marks it complete, updates progress |
| P2-T5 | **Progress Persistence** | Completed tasks remain checked after page refresh |
| P2-T6 | **Article Difficulty Tags** | Each article shows Easy/Medium/Hard badge |
| P2-T7 | **Daily Goal Progress** | Progress bar shows % of daily goal completed |

**Phase 2 Definition of Done:**
- [ ] All P2-T1 to P2-T7 tests pass
- [ ] Phase 1 tests still pass

---

## Phase 3: Progress Tracking (Week 2)

### 3.1 Progress Backend
- [ ] 3.1.1 Create `progress_logs` table
- [ ] 3.1.2 Add logging to all learning activities
- [ ] 3.1.3 Add `GET /api/progress/overview` - aggregate stats
- [ ] 3.1.4 Add `GET /api/progress/vocabulary` - vocabulary growth data

### 3.2 Streak System
- [ ] 3.2.1 Create `learning_streaks` table
- [ ] 3.2.2 Implement streak update logic
- [ ] 3.2.3 Add `GET /api/progress/streaks` endpoint

### 3.3 UI: Progress Dashboard
- [ ] 3.3.1 Add progress tab with key metrics
- [ ] 3.3.2 Add vocabulary growth chart (line chart)
- [ ] 3.3.3 Add streak display with fire icon 🔥
- [ ] 3.3.4 Add weekly summary view

### ✅ Phase 3 Acceptance Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| P3-T1 | **Activity Logging** | Reading an article creates a progress log entry |
| P3-T2 | **Vocabulary Count** | Dashboard shows "Words Learned: X" with correct count |
| P3-T3 | **Streak Start** | First activity of the day starts/continues streak |
| P3-T4 | **Streak Display** | Shows "🔥 X days" with current streak count |
| P3-T5 | **Streak Reset** | Missing a day resets streak to 0 (test with mock date) |
| P3-T6 | **Vocabulary Chart** | Line chart shows vocabulary growth over past 7 days |
| P3-T7 | **Weekly Summary** | Shows this week's: time spent, articles read, words learned |

**Phase 3 Definition of Done:**
- [ ] All P3-T1 to P3-T7 tests pass
- [ ] Phase 1 & 2 tests still pass

---

## Phase 4: Spaced Repetition System (Week 2-3)

### 4.1 SRS Backend
- [ ] 4.1.1 Create `vocabulary_reviews` table
- [ ] 4.1.2 Implement SM-2 algorithm in `internal/srs/sm2.go`
- [ ] 4.1.3 Add `GET /api/reviews/due` - get due vocabulary
- [ ] 4.1.4 Add `POST /api/reviews/answer` - update review schedule
- [ ] 4.1.5 Auto-add extracted vocabulary to review queue

### 4.2 UI: Review Mode
- [ ] 4.2.1 Add "Review Due" badge on dashboard
- [ ] 4.2.2 Create flashcard-style review interface
- [ ] 4.2.3 Add self-assessment buttons (Again/Hard/Good/Easy)
- [ ] 4.2.4 Show review statistics after session

### ✅ Phase 4 Acceptance Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| P4-T1 | **Vocabulary Auto-Add** | Extracting vocabulary auto-adds to review queue |
| P4-T2 | **Due Items Display** | Dashboard shows "X items due for review" |
| P4-T3 | **Flashcard Show/Hide** | Card shows word first, tap reveals meaning |
| P4-T4 | **Review - Easy** | Marking "Easy" schedules next review > 4 days later |
| P4-T5 | **Review - Again** | Marking "Again" schedules review for tomorrow |
| P4-T6 | **SM-2 Interval Growth** | After multiple "Good" reviews, interval increases |
| P4-T7 | **Review Stats** | After session shows: reviewed count, accuracy rate |
| P4-T8 | **No Due Items** | When nothing due, shows "All caught up! 🎉" |

**Phase 4 Definition of Done:**
- [ ] All P4-T1 to P4-T8 tests pass
- [ ] Phase 1, 2 & 3 tests still pass

---

## Phase 5: Intelligent Feedback (Week 3)

### 5.1 Feedback Backend
- [ ] 5.1.1 Create feedback chain in `internal/agent/feedback.go`
- [ ] 5.1.2 Implement answer evaluation for exercises
- [ ] 5.1.3 Generate personalized feedback with AI

### 5.2 Achievement System
- [ ] 5.2.1 Create `achievements` table
- [ ] 5.2.2 Define badge types and unlock criteria
- [ ] 5.2.3 Implement badge check on progress update
- [ ] 5.2.4 Add achievement notifications

### 5.3 UI: Feedback & Achievements
- [ ] 5.3.1 Show detailed feedback on exercise answers
- [ ] 5.3.2 Add achievement gallery/trophy case
- [ ] 5.3.3 Add celebration animation for new badges

### ✅ Phase 5 Acceptance Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| P5-T1 | **Correct Answer Feedback** | Shows ✓ and encouraging message |
| P5-T2 | **Wrong Answer Feedback** | Shows correct answer + explanation |
| P5-T3 | **7-Day Streak Badge** | After 7 consecutive days, earns "Week Warrior" badge |
| P5-T4 | **100 Words Badge** | Learning 100 words earns "Vocabulary Builder" badge |
| P5-T5 | **Badge Notification** | New badge shows celebration popup |
| P5-T6 | **Achievement Gallery** | Profile shows all earned badges with dates |
| P5-T7 | **Feedback Personalization** | Feedback mentions user's specific mistake pattern |

**Phase 5 Definition of Done:**
- [ ] All P5-T1 to P5-T7 tests pass
- [ ] Phase 1, 2, 3 & 4 tests still pass

---

## Phase 6: Interactive Exercises (Week 3-4)

### 6.1 Exercise Generation
- [ ] 6.1.1 Create exercise generator in `internal/agent/exercises.go`
- [ ] 6.1.2 Generate fill-in-the-blank exercises
- [ ] 6.1.3 Generate reading comprehension questions
- [ ] 6.1.4 Add `POST /api/exercises/generate` endpoint

### 6.2 Exercise Evaluation
- [ ] 6.2.1 Implement answer evaluation chain
- [ ] 6.2.2 Add `POST /api/exercises/evaluate` endpoint
- [ ] 6.2.3 Create `quiz_results` table and save results

### 6.3 UI: Exercise Mode
- [ ] 6.3.1 Add "Practice" tab/mode in UI
- [ ] 6.3.2 Create fill-in-the-blank UI component
- [ ] 6.3.3 Create multiple-choice comprehension UI
- [ ] 6.3.4 Show instant feedback with explanations
- [ ] 6.3.5 Display exercise session summary

### ✅ Phase 6 Acceptance Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| P6-T1 | **Exercise Generation** | After reading article, can generate 3 exercises |
| P6-T2 | **Fill-in-Blank Display** | Shows sentence with blank, accepts text input |
| P6-T3 | **Multiple Choice Display** | Shows question with 4 clickable options |
| P6-T4 | **Correct Answer Scoring** | Correct answer adds to score, shows ✓ |
| P6-T5 | **Wrong Answer Learning** | Wrong answer shows correct + explanation |
| P6-T6 | **Session Summary** | After exercises, shows: score, time, weak areas |
| P6-T7 | **Difficulty Match** | Exercises match user's assessed level |
| P6-T8 | **Exercise History** | Can view past exercise results in progress tab |

**Phase 6 Definition of Done:**
- [ ] All P6-T1 to P6-T8 tests pass
- [ ] All previous phase tests still pass (P1-P5)

---

## Final Validation

### Integration Tests
- [ ] Full user journey: Register → Assess → Plan → Read → Extract → Review → Exercise
- [ ] Multi-day simulation: streak, progress, SRS scheduling
- [ ] Performance: Dashboard loads in < 2 seconds

### Documentation
- [ ] Update README with new features
- [ ] Add user guide for each new feature
- [ ] Document API changes

---

## Test Execution Checklist

When completing each phase:
1. [ ] Run all acceptance tests for current phase
2. [ ] Run regression tests for all previous phases
3. [ ] Manual smoke test of main user flows
4. [ ] Update tasks.md to mark phase complete
5. [ ] Commit with message: `feat: complete Phase X - [description]`
