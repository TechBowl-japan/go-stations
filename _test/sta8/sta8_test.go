package sta8_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mattn/go-sqlite3"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

func TestStation8(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		Subject     string
		Description string
	}{
		"Subject is empty": {},
		"Description is empty": {
			Subject: "todo subject",
		},
		"Have not empty arguments": {
			Subject:     "todo subject",
			Description: "todo description",
		},
	}

	dbpath := "./temp_test.db"
	d, err := db.NewDB(dbpath)
	if err != nil {
		t.Error("DBの作成に失敗しました。", err)
		return
	}

	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Error("DBのクローズに失敗しました。", err)
			return
		}
		if err := os.Remove(dbpath); err != nil {
			t.Errorf("テスト用のDBファイルの削除に失敗しました: %v", err)
			return
		}
	})

	var sqlite3Err sqlite3.Error

	for name, tc := range testcases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			svc := service.NewTODOService(d)
			got, err := svc.CreateTODO(context.Background(), tc.Subject, tc.Description)
			if err != nil {
				if !errors.As(err, &sqlite3Err) {
					t.Errorf("期待していないエラーの Type です, got = %t, want = %+v", err, sqlite3Err)
				}
				return
			}

			want := &model.TODO{
				Subject:     tc.Subject,
				Description: tc.Description,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			if diff := cmp.Diff(got, want, cmpopts.EquateApproxTime(time.Second), cmpopts.IgnoreFields(model.TODO{}, "ID")); diff != "" {
				t.Error("期待していない値です\n", diff)
				return
			}
		})
	}
}
