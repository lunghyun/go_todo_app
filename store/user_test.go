package store

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lunghyun/go_todo_app/clock"
	"github.com/lunghyun/go_todo_app/entity"
)

func TestRepository_RegisterUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	c := clock.FixedClocker{}
	var wantID int64 = 20
	okUser := &entity.User{
		Name:     "okMello",
		Password: "password",
		Role:     "User",
		Created:  c.Now(),
		Modified: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	sql := `INSERT INTO user(name, password, role, created, modified) VALUES (?, ?, ?, ?, ?)`
	mock.ExpectExec(regexp.QuoteMeta(sql)).
		WithArgs(okUser.Name, okUser.Password, okUser.Role, okUser.Created, okUser.Modified).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}

	if err = r.RegisterUser(ctx, xdb, okUser); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
	if okUser.ID != entity.UserID(wantID) {
		t.Errorf("want ID %d, but got ID %d", wantID, okUser.ID)
	}
}

// TODO: RegisterUser: err case
func TestRepository_RegisterUser_DuplicateEntry(t *testing.T) {
	t.Parallel()
}

// TODO: RegisterUser: err case
func TestRepository_RegisterUser_Exec(t *testing.T) {
	t.Parallel()
}

// TODO: RegisterUser: err case
func TestRepository_RegisterUser_LastInsert(t *testing.T) {
	t.Parallel()
}
