package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/thedevrems/tuto_lua/server/internal/httpx"
	"github.com/thedevrems/tuto_lua/server/internal/models"
	"github.com/thedevrems/tuto_lua/server/internal/store"
)

// created is the uniform response for a successful content creation.
func created(w http.ResponseWriter, id string) {
	httpx.JSON(w, http.StatusCreated, map[string]string{"id": id})
}

// CreateCourse adds a new course (admin authoring).
func (h *AdminHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var c models.Course
	if !bind(w, r, &c) {
		return
	}
	id, err := h.store.CreateCourse(c)
	if errors.Is(err, store.ErrConflict) {
		httpx.Error(w, http.StatusConflict, "ce slug de cours est déjà utilisé")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "création du cours impossible")
		return
	}
	created(w, id)
}

// CreateChapter adds a chapter to the course named in the URL.
func (h *AdminHandler) CreateChapter(w http.ResponseWriter, r *http.Request) {
	var c models.Chapter
	if !bind(w, r, &c) {
		return
	}
	c.CourseID = chi.URLParam(r, "courseId")
	h.create(w, func() (string, error) { return h.store.CreateChapter(c) })
}

// CreateLesson adds a lesson to the chapter named in the URL.
func (h *AdminHandler) CreateLesson(w http.ResponseWriter, r *http.Request) {
	var l models.Lesson
	if !bind(w, r, &l) {
		return
	}
	l.ChapterID = chi.URLParam(r, "chapterId")
	h.create(w, func() (string, error) { return h.store.CreateLesson(l) })
}

// CreateExercise adds an exercise (with optional hints) to a chapter.
func (h *AdminHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	var e models.Exercise
	if !bind(w, r, &e) {
		return
	}
	e.ChapterID = chi.URLParam(r, "chapterId")
	h.create(w, func() (string, error) { return h.store.CreateExercise(e) })
}

// CreateTest adds an automated test to the exercise named in the URL.
func (h *AdminHandler) CreateTest(w http.ResponseWriter, r *http.Request) {
	var t models.ExerciseTest
	if !bind(w, r, &t) {
		return
	}
	t.ExerciseID = chi.URLParam(r, "exerciseId")
	h.create(w, func() (string, error) { return h.store.CreateTest(t) })
}

// create runs an insert closure and writes the standard 201/{id} or 500 response.
func (h *AdminHandler) create(w http.ResponseWriter, insert func() (string, error)) {
	id, err := insert()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "création impossible (élément parent introuvable ?)")
		return
	}
	created(w, id)
}

// mutate runs an update/delete closure and maps the outcome to a status:
// 404 (missing), 409 (slug conflict), 500 (other) or 204 (success).
func (h *AdminHandler) mutate(w http.ResponseWriter, op func() error) {
	switch err := op(); {
	case errors.Is(err, store.ErrNotFound):
		httpx.Error(w, http.StatusNotFound, "élément introuvable")
	case errors.Is(err, store.ErrConflict):
		httpx.Error(w, http.StatusConflict, "ce slug de cours est déjà utilisé")
	case err != nil:
		httpx.Error(w, http.StatusInternalServerError, "opération impossible")
	default:
		httpx.JSON(w, http.StatusNoContent, nil)
	}
}

// UpdateCourse edits an existing course.
func (h *AdminHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	var c models.Course
	if !bind(w, r, &c) {
		return
	}
	c.ID = chi.URLParam(r, "courseId")
	h.mutate(w, func() error { return h.store.UpdateCourse(c) })
}

// DeleteCourse removes a course and all its content.
func (h *AdminHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	h.mutate(w, func() error { return h.store.DeleteCourse(chi.URLParam(r, "courseId")) })
}

// UpdateChapter edits an existing chapter.
func (h *AdminHandler) UpdateChapter(w http.ResponseWriter, r *http.Request) {
	var c models.Chapter
	if !bind(w, r, &c) {
		return
	}
	c.ID = chi.URLParam(r, "chapterId")
	h.mutate(w, func() error { return h.store.UpdateChapter(c) })
}

// DeleteChapter removes a chapter and its lessons/exercises.
func (h *AdminHandler) DeleteChapter(w http.ResponseWriter, r *http.Request) {
	h.mutate(w, func() error { return h.store.DeleteChapter(chi.URLParam(r, "chapterId")) })
}

// UpdateLesson edits an existing lesson.
func (h *AdminHandler) UpdateLesson(w http.ResponseWriter, r *http.Request) {
	var l models.Lesson
	if !bind(w, r, &l) {
		return
	}
	l.ID = chi.URLParam(r, "lessonId")
	h.mutate(w, func() error { return h.store.UpdateLesson(l) })
}

// DeleteLesson removes a lesson.
func (h *AdminHandler) DeleteLesson(w http.ResponseWriter, r *http.Request) {
	h.mutate(w, func() error { return h.store.DeleteLesson(chi.URLParam(r, "lessonId")) })
}

// UpdateExercise edits an existing exercise.
func (h *AdminHandler) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	var e models.Exercise
	if !bind(w, r, &e) {
		return
	}
	e.ID = chi.URLParam(r, "exerciseId")
	h.mutate(w, func() error { return h.store.UpdateExercise(e) })
}

// DeleteExercise removes an exercise and its tests.
func (h *AdminHandler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	h.mutate(w, func() error { return h.store.DeleteExercise(chi.URLParam(r, "exerciseId")) })
}

// UpdateTest edits an existing automated test.
func (h *AdminHandler) UpdateTest(w http.ResponseWriter, r *http.Request) {
	var t models.ExerciseTest
	if !bind(w, r, &t) {
		return
	}
	t.ID = chi.URLParam(r, "testId")
	h.mutate(w, func() error { return h.store.UpdateTest(t) })
}

// DeleteTest removes an automated test.
func (h *AdminHandler) DeleteTest(w http.ResponseWriter, r *http.Request) {
	h.mutate(w, func() error { return h.store.DeleteTest(chi.URLParam(r, "testId")) })
}
