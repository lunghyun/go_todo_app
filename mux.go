package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/lunghyun/go_todo_app/handler"
	"github.com/lunghyun/go_todo_app/store"
)

func NewMux() http.Handler {
	//mux := http.NewServeMux()
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	v := validator.New()
	at := &handler.AddTask{Store: store.Tasks, Validator: v}
	mux.Post("/tasks", at.ServeHTTP)
	lt := &handler.ListTask{Store: store.Tasks}
	mux.Get("/tasks", lt.ServeHTTP)

	return mux
}
