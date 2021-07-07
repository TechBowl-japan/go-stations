package sta18_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

func TestStation18(t *testing.T) {
	t.Parallel()

	todos := []*model.TODO{
		{
			ID:      3,
			Subject: "todo subject 3",
		},
		{
			ID:      2,
			Subject: "todo subject 2",
		},
		{
			ID:      1,
			Subject: "todo subject 1",
		},
	}

	dbpath := "./temp_test.db"
	d, err := db.NewDB(dbpath)
	if err != nil {
		t.Error(" エラーが発生しました", err)
		return
	}

	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Error(" エラーが発生しました", err)
			return
		}
	})
	t.Cleanup(func() {
		if err := os.Remove(dbpath); err != nil {
			t.Error(" エラーが発生しました", err)
			return
		}
	})

	stmt, err := d.Prepare(`INSERT INTO todos(subject, description) VALUES(?, ?)`)
	if err != nil {
		t.Error(err)
		return
	}

	t.Cleanup(func() {
		if err := stmt.Close(); err != nil {
			t.Error("エラーが発生しました", err)
			return
		}
	})

	for _, todo := range []*model.TODO{todos[2], todos[1], todos[0]} {
		_, err = stmt.Exec(todo.Subject, todo.Description)
		if err != nil {
			t.Error("エラーが発生しました", err)
			return
		}
	}

	testcases := map[string]struct {
		IDs       []int64
		WantError error
	}{
		"Not found ID": {
			IDs:       []int64{4},
			WantError: &model.ErrNotFound{},
		},
		"One delete": {
			IDs:       []int64{1},
			WantError: nil,
		},
		"Multiple delete": {
			IDs:       []int64{2, 3},
			WantError: nil,
		},
		"Empty IDs": {
			IDs:       []int64{},
			WantError: nil,
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := service.NewTODOService(d).DeleteTODO(context.Background(), tc.IDs)

			switch tc.WantError {
			case nil:
				if err != nil {
					t.Error("エラーが発生しました", err)
					return
				}
			default:
				if !errors.As(err, &tc.WantError) {
					t.Errorf("期待していないエラーの型です, got = %+v, want = %+v", err, tc.WantError)
					return
				}
			}
		})
	}
}
