CREATE TABLE IF NOT EXISTS problems (
    path TEXT PRIMARY KEY,
    section TEXT NOT NULL,
    name TEXT NOT NULL UNIQUE,
    kyu INTEGER,
    tags TEXT,
    source TEXT,
    url TEXT,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    problem_path TEXT NOT NULL REFERENCES problems(path),
    date TEXT NOT NULL,
    passed INTEGER NOT NULL,
    duration_seconds INTEGER NOT NULL,
    started_at TEXT NOT NULL,
    finished_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS review_schedule (
    problem_path TEXT PRIMARY KEY REFERENCES problems(path),
    next_review TEXT NOT NULL,
    interval_days INTEGER NOT NULL DEFAULT 1,
    ease_factor REAL NOT NULL DEFAULT 2.5
);

CREATE INDEX IF NOT EXISTS idx_sessions_date ON sessions(date);
CREATE INDEX IF NOT EXISTS idx_sessions_problem ON sessions(problem_path);
CREATE INDEX IF NOT EXISTS idx_review_next ON review_schedule(next_review);
