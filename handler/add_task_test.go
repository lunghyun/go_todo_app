package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
	"github.com/lunghyun/go_todo_app/entity"
	"github.com/lunghyun/go_todo_app/testutil"
)

// TODO: AddTask: err case
func TestAddTask(t *testing.T) {
	type want struct {
		status  int
		resFile string
	}
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				resFile: "testdata/add_task/ok_res.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				resFile: "testdata/add_task/bad_res.json.golden",
			},
		},
	}
	for n, tt := range tests {
		//tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)
			moq := &AddTaskServiceMock{}
			moq.AddTaskFunc = func(ctx context.Context, title string) (*entity.Task, error) {
				if tt.want.status != http.StatusOK {
					return nil, errors.New("error from mock")
				}
				return &entity.Task{ID: 1}, nil
			}
			sut := AddTask{
				Service:   moq,
				Validator: validator.New()}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(
				t, resp, tt.want.status, testutil.LoadFile(t, tt.want.resFile),
			)
		})
	}
}
