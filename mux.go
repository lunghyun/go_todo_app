package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/lunghyun/go_todo_app/clock"
	"github.com/lunghyun/go_todo_app/config"
	"github.com/lunghyun/go_todo_app/handler"
	"github.com/lunghyun/go_todo_app/service"
	"github.com/lunghyun/go_todo_app/store"
)

// NewMux creates and returns the HTTP handler for the application along with a cleanup function used to release resources initialized for the mux.
// 
// NewMux configures the HTTP routes (including /health, POST /tasks and GET /tasks) and initializes persistent storage using the provided context and configuration.
// The returned cleanup function must be called to release storage and other resources when the server is shut down.
// If storage initialization fails, NewMux returns (nil, cleanup, err) where cleanup may still need to be invoked to free partial resources.
func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	r := store.Repository{Clocker: clock.RealClocker{}}
	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/tasks", at.ServeHTTP)
	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: &r},
	}
	mux.Get("/tasks", lt.ServeHTTP)

	return mux, cleanup, nil
}
