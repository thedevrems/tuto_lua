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
