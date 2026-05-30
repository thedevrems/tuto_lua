package database

import (
	"database/sql"
	"fmt"
)

// schema is the full relational model. It is idempotent (IF NOT EXISTS) so
// Migrate can run on every boot without external migration tooling.
const schema = `
CREATE TABLE IF NOT EXISTS users (
  id            TEXT PRIMARY KEY,
  username      TEXT NOT NULL UNIQUE COLLATE NOCASE,
  email         TEXT NOT NULL UNIQUE COLLATE NOCASE,
  password_hash TEXT NOT NULL,
  role          TEXT NOT NULL DEFAULT 'user',
  created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS courses (
  id          TEXT PRIMARY KEY,
  slug        TEXT NOT NULL UNIQUE,
  title       TEXT NOT NULL,
  summary     TEXT NOT NULL DEFAULT '',
  price_cents INTEGER NOT NULL DEFAULT 0,
  currency    TEXT NOT NULL DEFAULT 'eur',
  published   INTEGER NOT NULL DEFAULT 0,
  position    INTEGER NOT NULL DEFAULT 0,
  created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS chapters (
  id        TEXT PRIMARY KEY,
  course_id TEXT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  title     TEXT NOT NULL,
  summary   TEXT NOT NULL DEFAULT '',
  position  INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS lessons (
  id         TEXT PRIMARY KEY,
  chapter_id TEXT NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
  title      TEXT NOT NULL,
  content    TEXT NOT NULL DEFAULT '',
  position   INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS exercises (
  id         TEXT PRIMARY KEY,
  chapter_id TEXT NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
  title      TEXT NOT NULL,
  difficulty TEXT NOT NULL DEFAULT 'facile',
  statement  TEXT NOT NULL DEFAULT '',
  starter    TEXT NOT NULL DEFAULT '',
  solution   TEXT NOT NULL DEFAULT '',
  hints      TEXT NOT NULL DEFAULT '[]',
  position   INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS exercise_tests (
  id          TEXT PRIMARY KEY,
  exercise_id TEXT NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
  name        TEXT NOT NULL,
  code        TEXT NOT NULL,
  position    INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS enrollments (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  course_id  TEXT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  granted_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  source     TEXT NOT NULL DEFAULT 'purchase',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id, course_id)
);

CREATE TABLE IF NOT EXISTS progress (
  id          TEXT PRIMARY KEY,
  user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  exercise_id TEXT NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
  code        TEXT NOT NULL DEFAULT '',
  completed   INTEGER NOT NULL DEFAULT 0,
  updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id, exercise_id)
);

CREATE TABLE IF NOT EXISTS payments (
  id                TEXT PRIMARY KEY,
  user_id           TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  course_id         TEXT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  stripe_session_id TEXT NOT NULL DEFAULT '',
  amount_cents      INTEGER NOT NULL DEFAULT 0,
  currency          TEXT NOT NULL DEFAULT 'eur',
  status            TEXT NOT NULL DEFAULT 'pending',
  created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notifications (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title      TEXT NOT NULL,
  body       TEXT NOT NULL DEFAULT '',
  link       TEXT NOT NULL DEFAULT '',
  read       INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_chapters_course   ON chapters(course_id);
CREATE INDEX IF NOT EXISTS idx_lessons_chapter   ON lessons(chapter_id);
CREATE INDEX IF NOT EXISTS idx_exercises_chapter ON exercises(chapter_id);
CREATE INDEX IF NOT EXISTS idx_tests_exercise    ON exercise_tests(exercise_id);
CREATE INDEX IF NOT EXISTS idx_progress_user     ON progress(user_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_user  ON enrollments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_user     ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
`

// Migrate creates every table and index if they do not already exist.
func Migrate(db *sql.DB) error {
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}
	return nil
}
