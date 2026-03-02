package store

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/lunghyun/go_todo_app/clock"
	"github.com/lunghyun/go_todo_app/entity"
	"github.com/lunghyun/go_todo_app/testutil"
)

func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()
	// entity.Task를 작성하느 다른 테스트 케이스와 섞이면 테스트 실패
	// 따라서 트랜잭션을 적용 -> 테스트 케이스 내로 한정된 테이블 상태 만들기
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	// 해당 테스트 케이스가 끝나면 원래 상태로 되돌림
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}

	wants := prepareTasks(ctx, t, tx)
	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func TestRepository_ListTasks_Select(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled mock expectations: %v", err)
		}
		_ = db.Close()
	})

	sql := `SELECT id, title, status, created, modified FROM task`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WillReturnError(errors.New("select error"))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{}
	if _, err = r.ListTasks(ctx, xdb); err == nil {
		t.Error("want error, got nil")
	}
}

func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	okTask := &entity.Task{
		Title:    "ok task",
		Status:   entity.TaskStatusTodo,
		Created:  c.Now(),
		Modified: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled mock expectations: %v", err)
		}
		_ = db.Close()
	})

	sql := `INSERT INTO task (title, status, created, modified) VALUES (?, ?, ?, ?)`
	mock.ExpectExec(
		regexp.QuoteMeta(sql),
	).WithArgs(okTask.Title, okTask.Status, okTask.Created, okTask.Modified).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err = r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
	if okTask.ID != entity.TaskID(wantID) {
		t.Errorf("want ID %d, but got ID %d", wantID, okTask.ID)
	}
}

func TestRepository_AddTask_Exec(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	c := clock.FixedClocker{}

	errTask := &entity.Task{
		Title:    "err task",
		Status:   entity.TaskStatusTodo,
		Created:  c.Now(),
		Modified: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled mock expectations: %v", err)
		}
		_ = db.Close()
	})

	sql := `INSERT INTO task (title, status, created, modified) VALUES (?, ?, ?, ?)`
	mock.ExpectExec(
		regexp.QuoteMeta(sql),
	).WithArgs(errTask.Title, errTask.Status, errTask.Created, errTask.Modified).
		WillReturnError(errors.New("db task"))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err = r.AddTask(ctx, xdb, errTask); err == nil {
		t.Errorf("want error, but got nil")
	}
}

func TestRepository_AddTask_LastInsert(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	c := clock.FixedClocker{}

	errTask := &entity.Task{
		Title:    "err task",
		Status:   entity.TaskStatusTodo,
		Created:  c.Now(),
		Modified: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled mock expectations: %v", err)
		}
		_ = db.Close()
	})

	sql := `INSERT INTO task (title, status, created, modified) VALUES (?, ?, ?, ?)`
	mock.ExpectExec(
		regexp.QuoteMeta(sql),
	).WithArgs(errTask.Title, errTask.Status, errTask.Created, errTask.Modified).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err = r.AddTask(ctx, xdb, errTask); err == nil {
		t.Error("want error, but got nil")
	}
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) entity.Tasks {
	t.Helper()
	// 깨끗한 상태로 정리
	if _, err := con.ExecContext(ctx, "DELETE FROM task;"); err != nil {
		t.Logf("failed to initialize task: %v", err)
	}
	c := clock.FixedClocker{}

	wants := entity.Tasks{
		{
			Title:    "want task 1",
			Status:   entity.TaskStatusTodo,
			Created:  c.Now(),
			Modified: c.Now(),
		},
		{
			Title:    "want task 2",
			Status:   entity.TaskStatusTodo,
			Created:  c.Now(),
			Modified: c.Now(),
		},
		{
			Title:    "want task 3",
			Status:   entity.TaskStatusDone,
			Created:  c.Now(),
			Modified: c.Now(),
		},
	}
	result, err := con.ExecContext(ctx, "INSERT INTO task (title, status, created, modified) VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?);",
		wants[0].Title, wants[0].Status, wants[0].Created, wants[0].Modified,
		wants[1].Title, wants[1].Status, wants[1].Created, wants[1].Modified,
		wants[2].Title, wants[2].Status, wants[2].Created, wants[2].Modified,
	)
	if err != nil {
		t.Fatalf("unexpected error: result: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)
	return wants
}
