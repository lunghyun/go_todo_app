package handler

import (
	"log"
	"net/http"

	"github.com/lunghyun/go_todo_app/entity"
)

type ListTask struct {
	Service ListTasksService
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := lt.Service.ListTasks(ctx)
	if err != nil {
		log.Printf("ListTask failed: %+v ", err)
		RespondJSON(ctx, w, &ErrResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		}, http.StatusInternalServerError)
		return
	}
	res := make([]task, 0, len(tasks))
	for _, t := range tasks {
		res = append(res, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}
	RespondJSON(ctx, w, res, http.StatusOK)
}
