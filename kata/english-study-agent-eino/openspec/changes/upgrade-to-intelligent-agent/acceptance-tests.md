# Acceptance Tests: Intelligent English Learning Agent

> **Purpose**: Verify each phase meets requirements before proceeding to the next.
> **Strategy**: Each test should be executable manually or via automation.

---

## Phase 1: User Profile & Assessment

### P1-T1: New User Registration
**Precondition**: Clean database or new browser session
**Steps**:
1. Open the application
2. System should show onboarding wizard
3. Enter username: "test_user"
4. Select native language: "Chinese"
5. Set daily goal: 10 minutes
6. Select interests: ["Technology", "Business"]
7. Click "Start Learning"

**Expected**: User profile created, redirected to assessment or dashboard

---

### P1-T2: Profile Persistence
**Precondition**: P1-T1 completed
**Steps**:
1. Restart the backend server
2. Refresh the frontend page
3. Check profile in settings

**Expected**: Profile shows "test_user", Chinese, 10min goal, Technology & Business interests

---

### P1-T3: Assessment Start
**Precondition**: New user registration complete
**Steps**:
1. Click "Take Assessment" (or auto-start)
2. First question appears

**Expected**: Question is vocabulary-based, medium difficulty, has 4 answer options

---

### P1-T4: Adaptive Assessment
**Precondition**: Assessment started
**Steps**:
1. Answer first question correctly
2. Note second question difficulty
3. Answer second question incorrectly
4. Note third question difficulty

**Expected**: Q2 harder than Q1; Q3 easier than Q2

---

### P1-T5: Assessment Completion
**Precondition**: Answer all assessment questions
**Steps**:
1. Complete 5-10 questions
2. View results page

**Expected**: Shows CEFR level (e.g., "B1 - Intermediate"), breakdown by skill, recommendations

---

### P1-T6: Profile Update
**Precondition**: User has profile
**Steps**:
1. Go to Settings > Profile
2. Change daily goal to 15 minutes
3. Add "Medical" to interests
4. Save changes
5. Refresh page

**Expected**: Changes persist: 15min goal, interests include Medical

---

### P1-T7: Session Memory
**Precondition**: User logged in
**Steps**:
1. Close browser
2. Open browser again
3. Navigate to app

**Expected**: User is recognized, no registration prompt, dashboard loads with user data

---

## Phase 2: Learning Plan Generation

### P2-T1: Daily Plan Generation
**Precondition**: User has profile and assessment
**Steps**:
1. Open dashboard
2. View "Today's Tasks" section

**Expected**: 2-3 tasks listed with: title, estimated time, type (Read/Review/Practice)

---

### P2-T2: Plan Based on Level
**Precondition**: User assessed as A2 (Elementary)
**Steps**:
1. View recommended articles
2. Check difficulty badges

**Expected**: Most articles show "Easy" or "Medium" badges, no "Advanced"

---

### P2-T3: Plan Based on Interests
**Precondition**: User interests include "Technology"
**Steps**:
1. View recommended articles
2. Check article sources/categories

**Expected**: Technology articles appear before other categories

---

### P2-T4: Task Completion
**Precondition**: Daily plan visible
**Steps**:
1. Click checkbox on first task
2. Observe progress indicator

**Expected**: Task marked complete, progress bar increases

---

### P2-T5: Progress Persistence
**Precondition**: Task completed
**Steps**:
1. Refresh page
2. Check task status

**Expected**: Completed task still shows checkmark

---

### P2-T6: Article Difficulty Tags
**Precondition**: Article list visible
**Steps**:
1. View any 5 articles in list

**Expected**: Each article shows badge: "Easy" (green), "Medium" (yellow), or "Hard" (red)

---

### P2-T7: Daily Goal Progress
**Precondition**: Daily goal is 10 minutes
**Steps**:
1. Spend 5 minutes reading an article
2. Check progress indicator

**Expected**: Progress shows "5/10 min" or "50%" complete

---

## Phase 3: Progress Tracking

### P3-T1: Activity Logging
**Precondition**: User logged in
**Steps**:
1. Read an article for 2 minutes
2. Query database: `SELECT * FROM progress_logs ORDER BY id DESC LIMIT 1`

**Expected**: Log entry with activity_type='read', duration_seconds≈120

---

### P3-T2: Vocabulary Count
**Precondition**: User has extracted some vocabulary
**Steps**:
1. Extract 5 vocabulary words
2. View dashboard

**Expected**: "Words Learned" shows correct count including previous + 5

---

### P3-T3: Streak Start
**Precondition**: User has no activity today
**Steps**:
1. Complete any learning activity
2. Check streak display

**Expected**: Streak shows ≥1 day

---

### P3-T4: Streak Display
**Precondition**: User has 3-day streak
**Steps**:
1. View dashboard

**Expected**: Shows "🔥 3 days" prominently

---

### P3-T5: Streak Reset
**Precondition**: User has streak, simulate missing a day
**Steps**:
1. Modify last_activity_date in DB to 2 days ago
2. Perform activity today
3. Check streak

**Expected**: Streak resets to 1

---

### P3-T6: Vocabulary Chart
**Precondition**: User has vocabulary data over multiple days
**Steps**:
1. View Progress tab
2. Locate vocabulary chart

**Expected**: Line chart shows cumulative words over past 7 days

---

### P3-T7: Weekly Summary
**Precondition**: User has activity this week
**Steps**:
1. View Progress > Weekly Summary

**Expected**: Shows: total time, articles read, words learned, comparison to goal

---

## Phase 4: Spaced Repetition System

### P4-T1: Vocabulary Auto-Add
**Precondition**: Review queue is empty
**Steps**:
1. Read an article
2. Click "Extract Vocabulary"
3. Check review queue

**Expected**: Extracted words appear in review queue with next_review = tomorrow

---

### P4-T2: Due Items Display
**Precondition**: Items due for review today
**Steps**:
1. View dashboard

**Expected**: Badge shows "X items due" or notification

---

### P4-T3: Flashcard Show/Hide
**Precondition**: Start review session
**Steps**:
1. View flashcard (shows word)
2. Click "Show Answer"

**Expected**: Card reveals meaning, context sentence, and rating buttons

---

### P4-T4: Review - Easy
**Precondition**: Reviewing a card
**Steps**:
1. Click "Easy" button
2. Check next_review_at in database

**Expected**: Next review scheduled ≥4 days from now

---

### P4-T5: Review - Again
**Precondition**: Reviewing a card
**Steps**:
1. Click "Again" button
2. Check next_review_at in database

**Expected**: Next review scheduled for tomorrow

---

### P4-T6: SM-2 Interval Growth
**Precondition**: Card has been reviewed "Good" 3 times
**Steps**:
1. Review card again, mark "Good"
2. Check interval_days in database

**Expected**: Interval > previous interval (e.g., 1 → 6 → 15 → 30+ days)

---

### P4-T7: Review Stats
**Precondition**: Complete review session of 5 items
**Steps**:
1. Complete all reviews
2. View session summary

**Expected**: Shows: "5 reviewed", accuracy rate (e.g., "80%"), next review preview

---

### P4-T8: No Due Items
**Precondition**: All reviews completed
**Steps**:
1. Try to start review session

**Expected**: Shows "All caught up! 🎉" or "No reviews due"

---

## Phase 5: Intelligent Feedback

### P5-T1: Correct Answer Feedback
**Precondition**: Exercise in progress
**Steps**:
1. Answer question correctly

**Expected**: Shows ✓, encouraging message (e.g., "Great job!")

---

### P5-T2: Wrong Answer Feedback
**Precondition**: Exercise in progress
**Steps**:
1. Answer question incorrectly

**Expected**: Shows correct answer, explanation of why it's correct

---

### P5-T3: 7-Day Streak Badge
**Precondition**: User has 6-day streak
**Steps**:
1. Complete activity on day 7

**Expected**: "Week Warrior" badge awarded, celebration popup

---

### P5-T4: 100 Words Badge
**Precondition**: User has 99 words learned
**Steps**:
1. Extract 1 more vocabulary word

**Expected**: "Vocabulary Builder" badge awarded

---

### P5-T5: Badge Notification
**Precondition**: Badge earned
**Steps**:
1. Observe UI after earning badge

**Expected**: Popup or notification celebrates the achievement

---

### P5-T6: Achievement Gallery
**Precondition**: User has earned badges
**Steps**:
1. Go to Profile > Achievements

**Expected**: Shows all earned badges with name, description, date earned

---

### P5-T7: Feedback Personalization
**Precondition**: User has made same mistake type 3+ times
**Steps**:
1. Make the same mistake again

**Expected**: Feedback mentions the pattern (e.g., "Watch out for article usage - this is a common challenge!")

---

## Phase 6: Interactive Exercises

### P6-T1: Exercise Generation
**Precondition**: Article read
**Steps**:
1. Click "Generate Exercises"

**Expected**: 3 exercises generated based on article content

---

### P6-T2: Fill-in-Blank Display
**Precondition**: Fill-in-blank exercise loaded
**Steps**:
1. View exercise

**Expected**: Sentence with blank (____), text input field

---

### P6-T3: Multiple Choice Display
**Precondition**: Multiple choice exercise loaded
**Steps**:
1. View exercise

**Expected**: Question with 4 clickable option buttons

---

### P6-T4: Correct Answer Scoring
**Precondition**: Exercise in progress
**Steps**:
1. Select correct answer

**Expected**: Score increases, shows ✓

---

### P6-T5: Wrong Answer Learning
**Precondition**: Exercise in progress
**Steps**:
1. Select wrong answer

**Expected**: Shows correct answer + brief explanation

---

### P6-T6: Session Summary
**Precondition**: Complete exercise session
**Steps**:
1. Finish all exercises

**Expected**: Summary shows: score (e.g., 4/5), time taken, areas to improve

---

### P6-T7: Difficulty Match
**Precondition**: User assessed as B1
**Steps**:
1. Generate exercises
2. Observe difficulty

**Expected**: Exercises are intermediate level, not too easy or too hard

---

### P6-T8: Exercise History
**Precondition**: Completed exercises
**Steps**:
1. Go to Progress > Exercise History

**Expected**: List of past exercises with scores and dates

---

## Regression Test Suite

After each phase, verify these core features still work:

| Feature | Quick Test |
|---------|------------|
| RSS Feed Loading | Articles list loads from configured feeds |
| Article Reading | Can open and read full article |
| AI Explain | "Explain" button returns AI explanation |
| AI Translate | "Translate" button returns Chinese translation |
| Vocabulary Extract | "Extract" returns phrases and sentences |
| TTS | "Read" button speaks text aloud |
| Custom RSS | Can add/edit/delete custom feeds |
| Settings | Can change language (EN/ZH) |

