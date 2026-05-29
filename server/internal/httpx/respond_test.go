package httpx

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSONWritesStatusAndBody(t *testing.T) {
	rec := httptest.NewRecorder()
	JSON(rec, 201, map[string]int{"n": 7})
	if rec.Code != 201 {
		t.Fatalf("status = %d, want 201", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("content-type = %q", ct)
	}
	if !strings.Contains(rec.Body.String(), `"n":7`) {
		t.Fatalf("body = %q", rec.Body.String())
	}
}

func TestErrorEnvelope(t *testing.T) {
	rec := httptest.NewRecorder()
	Error(rec, 400, "boom")
	if !strings.Contains(rec.Body.String(), `"error":"boom"`) {
		t.Fatalf("body = %q", rec.Body.String())
	}
}

func TestDecodeRejectsUnknownFields(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"surprise":1}`))
	var dst struct {
		Known string `json:"known"`
	}
	if err := Decode(rec, req, &dst); err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestDecodeRejectsTrailingContent(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"known":"a"}{}`))
	var dst struct {
		Known string `json:"known"`
	}
	if err := Decode(rec, req, &dst); err == nil {
		t.Fatal("expected error for trailing JSON")
	}
}

func TestDecodeAcceptsValidBody(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"known":"hello"}`))
	var dst struct {
		Known string `json:"known"`
	}
	if err := Decode(rec, req, &dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst.Known != "hello" {
		t.Fatalf("Known = %q", dst.Known)
	}
}
