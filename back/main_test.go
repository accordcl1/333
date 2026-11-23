package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestTasksAPI(t *testing.T) {
	// Используем временный файл для теста, чтобы не трогать /data
	dataFile = filepath.Join(t.TempDir(), "test_tasks.json")
	tasks = nil

	// 1. Проверяем создание задачи (POST)
	body := bytes.NewBufferString(`{"text":"Deploy to prod"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	w := httptest.NewRecorder()
	tasksHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	// 2. Проверяем получение списка (GET)
	reqGet := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	wGet := httptest.NewRecorder()
	tasksHandler(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", wGet.Code)
	}
}

