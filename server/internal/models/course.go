package models

import "time"

// Course is a purchasable learning track that groups chapters.
type Course struct {
	ID         string    `json:"id"`
	Slug       string    `json:"slug"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	PriceCents int       `json:"priceCents"`
	Currency   string    `json:"currency"`
	Published  bool      `json:"published"`
	Position   int       `json:"position"`
	CreatedAt  time.Time `json:"createdAt"`
	// Chapters is populated only when a full course tree is requested.
	Chapters []Chapter `json:"chapters,omitempty"`
}

// Chapter groups lessons and exercises inside a course.
type Chapter struct {
	ID        string     `json:"id"`
	CourseID  string     `json:"courseId"`
	Title     string     `json:"title"`
	Summary   string     `json:"summary"`
	Position  int        `json:"position"`
	Lessons   []Lesson   `json:"lessons,omitempty"`
	Exercises []Exercise `json:"exercises,omitempty"`
}

// Lesson is reading material rendered as markdown.
type Lesson struct {
	ID        string `json:"id"`
	ChapterID string `json:"chapterId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Position  int    `json:"position"`
}

// Exercise is an interactive coding task with an editor and tests.
type Exercise struct {
	ID         string         `json:"id"`
	ChapterID  string         `json:"chapterId"`
	Title      string         `json:"title"`
	Difficulty string         `json:"difficulty"`
	Statement  string         `json:"statement"`
	Starter    string         `json:"starter"`
	// Solution is omitted from listings and only sent on demand.
	Solution string         `json:"solution,omitempty"`
	Position int            `json:"position"`
	Tests    []ExerciseTest `json:"tests,omitempty"`
}

// ExerciseTest is one automated assertion run against the student's code.
type ExerciseTest struct {
	ID         string `json:"id"`
	ExerciseID string `json:"exerciseId"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Position   int    `json:"position"`
}
