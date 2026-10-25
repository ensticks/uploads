package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type bindReq struct {
	Name string `json:"name"`
}

func TestBindJSON_Success(t *testing.T) {
	body := bytes.NewBufferString(`{"name":"alice"}`)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/bind", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	var dst bindReq
	ok := BindJSON(w, req, &dst)

	if !ok {
		t.Fatalf("expected BindJSON to succeed, got status %d: %s", w.Code, w.Body.String())
	}
	if dst.Name != "alice" {
		t.Fatalf("expected name=alice, got %q", dst.Name)
	}
}

func TestBindJSON_Failure(t *testing.T) {
	body := bytes.NewBufferString(`not-json`)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/bind", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	var dst bindReq
	ok := BindJSON(w, req, &dst)

	if ok {
		t.Fatal("expected BindJSON to fail on invalid JSON")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d (%s)", w.Code, w.Body.String())
	}

	var env struct {
		ErrorCode string `json:"error_code"`
	}
	if err := json.NewDecoder(w.Body).Decode(&env); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if env.ErrorCode != "BAD_REQUEST" {
		t.Fatalf("expected BAD_REQUEST, got %q", env.ErrorCode)
	}
}
