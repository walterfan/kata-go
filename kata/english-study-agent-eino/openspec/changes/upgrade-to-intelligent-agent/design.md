# Design: Intelligent English Learning Agent

## Context

当前系统功能对比 | 目标状态：

| 当前 (Language Tool) | 目标 (AI Agent) |
|---------------------|-----------------|
| 被动响应用户请求 | 主动推荐学习内容 |
| 无用户概念 | 个性化用户档案 |
| 无水平评估 | 入学测试 + 定期评估 |
| 无学习计划 | 自适应学习路径 |
| 无进度跟踪 | 学习曲线可视化 |
| 无记忆巩固 | 间隔重复系统 |
| 单向输出 | 互动练习与反馈 |

Target Users: Non-native English speakers, especially software engineers, seeking efficient daily learning (10-15 min/day).

## Goals / Non-Goals

### Goals
1. **Understand the User**: Create comprehensive learner profile
2. **Personalize Learning**: Adaptive content and difficulty
3. **Track Progress**: Visualize improvement over time
4. **Optimize Retention**: Spaced repetition for vocabulary
5. **Provide Feedback**: Real-time corrections and encouragement
6. **Stay Simple**: MVP in 2 weeks, iterate based on feedback

### Non-Goals
- Full SRS implementation like Anki (use simplified SM-2)
- Speech recognition for pronunciation (text-based first)
- Gamification beyond basic badges (keep it professional)
- Mobile app (web-responsive only for now)

## Architecture

### Agent Intelligence Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    User Interface Layer                      │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌───────────┐│
│  │ Dashboard  │ │  Learning  │ │  Exercises │ │  Review   ││
│  │ (Progress) │ │  (Read)    │ │  (Quiz)    │ │  (SRS)    ││
│  └────────────┘ └────────────┘ └────────────┘ └───────────┘│
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                    Agent Brain (Eino)                        │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                   Planner Chain                       │  │
│  │  • Assess User Level                                  │  │
│  │  • Generate Learning Goals                            │  │
│  │  • Select Content for Today                           │  │
│  │  • Create Practice Exercises                          │  │
│  └──────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                   Feedback Chain                      │  │
│  │  • Evaluate User Answers                              │  │
│  │  • Generate Encouraging Feedback                      │  │
│  │  • Identify Weak Areas                                │  │
│  │  • Suggest Next Steps                                 │  │
│  └──────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                   Analysis Chain                      │  │
│  │  • Explain Articles (existing)                        │  │
│  │  • Extract Vocabulary (existing)                      │  │
│  │  • Simplify Text (existing)                           │  │
│  │  • Translate (existing)                               │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                    Data Layer (SQLite)                       │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐   │
│  │  Users   │ │ Progress │ │ Reviews  │ │ Achievements │   │
│  │ Profiles │ │   Logs   │ │  (SRS)   │ │    Badges    │   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### New Database Schema

```sql
-- User Management
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_profiles (
    user_id INTEGER PRIMARY KEY REFERENCES users(id),
    native_language TEXT DEFAULT 'Chinese',
    target_level TEXT DEFAULT 'Intermediate', -- Beginner/Intermediate/Advanced
    daily_goal_minutes INTEGER DEFAULT 10,
    interests TEXT, -- JSON array: ["technology", "medicine", "business"]
    vocabulary_size_estimate INTEGER,
    reading_level TEXT, -- CEFR: A1-C2
    assessed_at DATETIME,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Learning Goals
CREATE TABLE learning_goals (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    goal_type TEXT, -- 'daily', 'weekly', 'monthly'
    description TEXT,
    target_value INTEGER,
    current_value INTEGER DEFAULT 0,
    start_date DATE,
    end_date DATE,
    status TEXT DEFAULT 'active', -- 'active', 'completed', 'failed'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Progress Tracking
CREATE TABLE progress_logs (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    activity_type TEXT, -- 'read', 'quiz', 'review', 'exercise'
    activity_data TEXT, -- JSON with details
    duration_seconds INTEGER,
    score REAL, -- 0-100 for quizzes
    logged_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Spaced Repetition
CREATE TABLE vocabulary_reviews (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    word TEXT NOT NULL,
    context TEXT, -- Original sentence
    translation TEXT,
    ease_factor REAL DEFAULT 2.5, -- SM-2 algorithm
    interval_days INTEGER DEFAULT 1,
    repetitions INTEGER DEFAULT 0,
    next_review_at DATETIME,
    last_reviewed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Quiz & Exercise Results
CREATE TABLE quiz_results (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    quiz_type TEXT, -- 'vocabulary', 'comprehension', 'grammar'
    question TEXT,
    user_answer TEXT,
    correct_answer TEXT,
    is_correct BOOLEAN,
    feedback TEXT,
    completed_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Achievements
CREATE TABLE achievements (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    badge_type TEXT, -- 'streak_7', 'words_100', 'articles_10'
    badge_name TEXT,
    description TEXT,
    earned_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Learning Streaks
CREATE TABLE learning_streaks (
    user_id INTEGER PRIMARY KEY REFERENCES users(id),
    current_streak INTEGER DEFAULT 0,
    longest_streak INTEGER DEFAULT 0,
    last_activity_date DATE
);
```

### New API Endpoints

```
# User Management
POST   /api/users/register
POST   /api/users/login
GET    /api/users/profile
PUT    /api/users/profile

# Assessment
POST   /api/assessment/start       # Start entry-level test
POST   /api/assessment/answer      # Submit answer
GET    /api/assessment/result      # Get assessment result

# Learning Plan
GET    /api/plan/today             # Get today's learning tasks
GET    /api/plan/weekly            # Get weekly goals
POST   /api/plan/complete-task     # Mark task complete

# Progress
GET    /api/progress/overview      # Dashboard stats
GET    /api/progress/vocabulary    # Vocabulary growth chart
GET    /api/progress/streaks       # Streak information

# Reviews (SRS)
GET    /api/reviews/due            # Get due reviews
POST   /api/reviews/answer         # Submit review answer
GET    /api/reviews/stats          # Review statistics

# Exercises
POST   /api/exercises/generate     # Generate exercises from article
POST   /api/exercises/evaluate     # Evaluate user answers
```

### Agent Prompts (New Chains)

#### Assessment Chain
```
You are an English proficiency assessor. Based on the user's answers,
determine their:
1. Vocabulary Level (A1-C2)
2. Reading Comprehension Level
3. Grammar Proficiency

Provide specific areas for improvement.
Output JSON: {level, vocabulary_score, reading_score, grammar_score, recommendations}
```

#### Planner Chain
```
You are an English learning planner for {user_name}.
User Profile: {profile_json}
Recent Progress: {progress_json}

Create today's learning plan:
1. Review due vocabulary (if any)
2. One reading exercise at their level
3. One practice exercise targeting weak areas

Output JSON: {tasks: [{type, content, estimated_minutes}]}
```

#### Feedback Chain
```
You are an encouraging English tutor.
The user answered: {user_answer}
Correct answer: {correct_answer}
Context: {context}

Provide:
1. Whether they were correct
2. Clear explanation if wrong
3. Encouragement and tip
4. Related vocabulary to learn

Be supportive but accurate. Use simple language.
```

## Decisions

### User Authentication
- **Decision**: Simple username-based local auth (no password for MVP)
- **Rationale**: Keep friction low; this is a personal learning tool
- **Future**: Add OAuth if multi-user needed

### Assessment Method
- **Decision**: Use AI-generated adaptive questions (5-10 questions)
- **Rationale**: More flexible than static tests; adjusts to responses
- **Implementation**: Start with medium difficulty, adjust based on answers

### Spaced Repetition
- **Decision**: Simplified SM-2 algorithm
- **Rationale**: Well-proven, simple to implement
- **Parameters**: Initial interval=1 day, ease_factor starts at 2.5

### Progress Visualization
- **Decision**: Use Streamlit's built-in charts (line, bar, metrics)
- **Rationale**: Keep it simple; no need for external charting libs

## Risks / Trade-offs

### Risk: Complexity Creep
- **Mitigation**: Strict phase-based delivery; ship Phase 1 before Phase 2
- **Validation**: Each phase must be usable standalone

### Risk: AI Costs Increase
- **Mitigation**: Cache aggressively; limit assessment frequency
- **Budget**: Target <$5/user/month at active usage

### Risk: User Drop-off
- **Mitigation**: Keep daily commitment low (10 min); add streaks for motivation
- **Metric**: Track 7-day retention

## Migration Plan

1. **Phase 0**: Create migration script for existing DB
   - Add `users` table with default user (id=1)
   - Add `user_id=1` to existing `learning_items`
   - No data loss

2. **Each Phase**: Additive changes only; old features continue working

## Open Questions

1. **Profile Data**: Should we ask detailed questions upfront or infer from usage?
   - **Lean**: Start with 3 simple questions (native language, goal, time), infer rest

2. **Offline Support**: Should reviews work offline?
   - **Decision**: Web-only for MVP; defer offline

3. **Multi-language UI**: Already have EN/ZH toggle; sufficient for now

