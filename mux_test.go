package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lunghyun/go_todo_app/config"
)

func TestNewMux(t *testing.T) {
	// 1. 서버 셋업
	ctx := context.Background()
	if _, defined := os.LookupEnv("CI"); defined {
		t.Setenv("TODO_DB_PORT", "3306")
	}
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	sut, cleanup, err := NewMux(ctx, cfg)
	t.Cleanup(cleanup)
	if err != nil {
		t.Fatalf("failed to create mux: %v", err)
	}

	// 2. 준비된 서버에서 요청, 응답 셋업
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	sut.ServeHTTP(w, r)

	res := w.Result()
	t.Cleanup(func() { _ = res.Body.Close() })

	if res.StatusCode != http.StatusOK {
		t.Errorf("want status code 200; but %v", res.StatusCode)
	}

	got, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := `{"status":"ok"}`
	if string(got) != want {
		t.Errorf("want %q; but %q", want, got)
	}
}
