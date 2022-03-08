package sta15_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

func TestStation15(t *testing.T) {
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
		t.Errorf("データベースの作成に失敗しました: %v", err)
		return
	}

	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Errorf("データベースのクローズに失敗しました: %v", err)
			return
		}
		if err := os.Remove(dbpath); err != nil {
			t.Errorf("テスト用のデータベースの削除に失敗しました: %v", err)
			return
		}
	})

	stmt, err := d.Prepare(`INSERT INTO todos(subject, description) VALUES(?, ?)`)
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
		if _, err := stmt.Exec(todo.Subject, todo.Description); err != nil {
			t.Errorf("データベースのステートメントの実行に失敗しました: %v", err)
			return
		}
	}

	testcases := map[string]struct {
		PrevID int64
		Size   int64
		TODOs  []*model.TODO
	}{
		"Zero read": {
			PrevID: 0,
			Size:   0,
			TODOs:  todos[3:],
		},
		"All read": {
			PrevID: 0,
			Size:   5,
			TODOs:  todos,
		},
		"One read": {
			PrevID: 0,
			Size:   1,
			TODOs:  todos[:1],
		},
		"All read with prev id = 3": {
			PrevID: 3,
			Size:   5,
			TODOs:  todos[1:],
		},
		"All read with prev id = 1": {
			PrevID: 1,
			Size:   5,
			TODOs:  todos[3:],
		},
		"One read with prev id = 3": {
			PrevID: 3,
			Size:   1,
			TODOs:  todos[1:2],
		},
		"One read with prev id = 1": {
			PrevID: 1,
			Size:   1,
			TODOs:  todos[3:],
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			svc := service.NewTODOService(d)
			ret, err := svc.ReadTODO(context.Background(), tc.PrevID, tc.Size)
			if err != nil {
				t.Errorf("ReadTODOに失敗しました: %v", err)
				return
			}
			if diff := cmp.Diff(ret, tc.TODOs, cmpopts.IgnoreFields(model.TODO{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Error("期待していない値です\n", diff)
				return
			}
		})
	}
}
