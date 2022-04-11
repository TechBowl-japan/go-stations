package sta8_test

import (
	"os"
	"testing"

	"github.com/TechBowl-japan/go-stations/db"
)

func TestStation8(t *testing.T) {
	t.Parallel()

	// testcases := map[string]struct {
	// 	Subject     string
	// 	Description string
	// 	WantError   error
	// }{
	// 	"Subject is empty": {
	// 		WantError: sqlite3.ErrConstraint,
	// 	},
	// 	"Description is empty": {
	// 		Subject: "todo subject",
	// 	},
	// 	"Have not empty arguments": {
	// 		Subject:     "todo subject",
	// 		Description: "todo description",
	// 	},
	// }

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

	// for name, tc := range testcases {
	// 	name := name
	// 	tc := tc
	// 	t.Run(name, func(t *testing.T) {
	// 		t.Parallel()

	// 		svc := service.NewTODOService(d)
	// 		got, err := svc.CreateTODO(context.Background(), tc.Subject, tc.Description)
	// 		switch tc.WantError {
	// 		case nil:
	// 			if err != nil {
	// 				t.Error("エラーが発生しました", err)
	// 				return
	// 			}
	// 		default:
	// 			if !errors.As(err, &tc.WantError) {
	// 				t.Errorf("期待していないエラーの Type です, got = %t, want = %+v", err, tc.WantError)
	// 			}
	// 			return
	// 		}

	// 		want := &model.TODO{
	// 			Subject:     tc.Subject,
	// 			Description: tc.Description,
	// 			CreatedAt:   time.Now(),
	// 			UpdatedAt:   time.Now(),
	// 		}
	// 		if diff := cmp.Diff(got, want, cmpopts.EquateApproxTime(time.Second), cmpopts.IgnoreFields(model.TODO{}, "ID")); diff != "" {
	// 			t.Error("期待していない値です\n", diff)
	// 			return
	// 		}
	// 	})
	// }
}
