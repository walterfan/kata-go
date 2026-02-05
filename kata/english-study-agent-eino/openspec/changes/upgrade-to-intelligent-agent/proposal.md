# Change: Upgrade to Intelligent English Learning Agent

## Why
当前系统是一个优秀的**语言分析工具**，但缺少真正AI Agent的核心特征：它不了解用户、不制定学习计划、没有测评和反馈、无法跟踪进度。真正的智能学习助手应该像私人教师一样，了解学生水平、制定个性化计划、评估学习效果、持续调整策略。

The current system is a good **language analysis tool**, but lacks the core characteristics of a true AI Agent: it doesn't understand users, doesn't create learning plans, has no assessments or feedback, and cannot track progress. A true intelligent learning assistant should act like a private tutor—understanding student levels, creating personalized plans, evaluating learning outcomes, and continuously adjusting strategies.

## What Changes

### Phase 1: User Profile & Assessment (用户档案与测评)
- Add user registration and profile management
- Add entry-level assessment (vocabulary, reading, grammar)
- Store user's learning history and preferences
- Track user's strengths and weaknesses

### Phase 2: Personalized Learning Plan (个性化学习计划)
- Generate daily/weekly learning goals based on user level
- Create adaptive learning paths targeting weak areas
- Recommend content based on user interests and level
- Balance reading, vocabulary, and comprehension exercises

### Phase 3: Progress Tracking & Analytics (进度跟踪与分析)
- Track vocabulary acquisition (words learned, retention rate)
- Monitor reading comprehension progress
- Visualize learning curve and achievements
- Generate weekly/monthly progress reports

### Phase 4: Intelligent Feedback System (智能反馈系统)
- Provide real-time feedback on exercises
- Identify common mistakes and suggest corrections
- Celebrate achievements with milestone badges
- Adaptive difficulty based on performance

### Phase 5: Spaced Repetition System (间隔重复系统)
- Implement SM-2 algorithm for vocabulary review
- Schedule reviews based on memory decay curve
- Integrate review reminders into daily workflow
- Track long-term retention rates

### Phase 6: Interactive Learning (互动学习)
- Add fill-in-the-blank exercises
- Add reading comprehension quizzes
- Add pronunciation practice (with TTS feedback)
- Add writing exercises with AI correction

## Impact

### Affected Specs
- `ai-assistant` - MODIFIED: Add personalization and assessment capabilities
- `user-profile` - NEW: User management and preferences
- `learning-plan` - NEW: Adaptive learning path generation  
- `progress-tracking` - NEW: Analytics and progress visualization
- `feedback-system` - NEW: Interactive feedback and corrections
- `spaced-repetition` - NEW: Memory optimization system
- `exercises` - NEW: Interactive practice modules

### Affected Code
- `internal/storage/db.go` - New tables for users, progress, reviews
- `internal/agent/` - New chains for assessment, planning, feedback
- `internal/api/server.go` - New endpoints for user management
- `web/app.py` - New dashboard, progress charts, exercise UI

### Database Schema Changes (**BREAKING**)
- New tables: `users`, `user_profiles`, `learning_goals`, `progress_logs`, `vocabulary_reviews`, `quiz_results`, `achievements`
- Modified tables: `learning_items` - add `user_id`, `mastery_level`, `next_review_at`

### Migration
- Existing data will be migrated to default user
- No data loss expected

