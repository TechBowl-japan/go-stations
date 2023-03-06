package sta18_test

import (
	"context"
	"os"
	"reflect"
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

	dbPath := "./temp_test.db"
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		t.Errorf("データベースの作成に失敗しました: %v", err)
		return
	}

	t.Cleanup(func() {
		if err := todoDB.Close(); err != nil {
			t.Errorf("データベースのクローズに失敗しました: %v", err)
			return
		}
		if err := os.Remove(dbPath); err != nil {
			t.Errorf("テスト用のDBファイルの削除に失敗しました: %v", err)
			return
		}
	})

	stmt, err := todoDB.Prepare(`INSERT INTO todos(subject, description) VALUES(?, ?)`)
	if err != nil {
		t.Errorf("データベースのステートメントの作成に失敗しました: %v", err)
		return
	}

	t.Cleanup(func() {
		if err := stmt.Close(); err != nil {
			t.Errorf("データベースのステートメントのクローズに失敗しました: %v", err)
			return
		}
	})

	for _, todo := range []*model.TODO{todos[2], todos[1], todos[0]} {
		if _, err = stmt.Exec(todo.Subject, todo.Description); err != nil {
			t.Errorf("データベースのステートメントの実行に失敗しました: %v", err)
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
	}

	for name, tc := range testcases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := service.NewTODOService(todoDB).DeleteTODO(context.Background(), tc.IDs)

			switch tc.WantError {
			case nil:
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
					return
				}
			default:
				if reflect.TypeOf(err) != reflect.TypeOf(tc.WantError) {
					t.Errorf("期待していないエラーの型です, got = %+v, want = %+v", err, tc.WantError)
					return
				}
			}
		})
	}
}
