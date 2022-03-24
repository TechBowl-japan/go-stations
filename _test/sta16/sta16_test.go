package sta16_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func TestStation16(t *testing.T) {
	dbPath := "./temp_test.db"
	if err := os.Setenv("DB_PATH", dbPath); err != nil {
		t.Errorf("dbPathのセットに失敗しました。%v", err)
		return
	}

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

	stmt, err := todoDB.Prepare(`INSERT INTO todos(subject) VALUES(?)`)
	if err != nil {
		t.Errorf("ステートメントの作成に失敗しました: %v", err)
		return
	}

	t.Cleanup(func() {
		if err := stmt.Close(); err != nil {
			t.Errorf("ステートメントのクローズに失敗しました: %v", err)
			return
		}
	})

	for i := 0; i < 3; i++ {
		if _, err := stmt.Exec(fmt.Sprintf("todo subject %d", i+1)); err != nil {
			t.Errorf("todoの追加に失敗しました: %v", err)
			return
		}
	}

	r := router.NewRouter(todoDB)
	srv := httptest.NewServer(r)
	defer srv.Close()

	testcases := map[string]struct {
		PrevID   int
		Size     int
		TODOsLen int
	}{
		"Default read": {
			PrevID:   0,
			Size:     0,
			TODOsLen: 3,
		},
		"All read": {
			PrevID:   0,
			Size:     5,
			TODOsLen: 3,
		},
		"One read": {
			PrevID:   0,
			Size:     1,
			TODOsLen: 1,
		},
		"All read with prev id = 3": {
			PrevID:   3,
			Size:     5,
			TODOsLen: 2,
		},
		"All read with prev id = 1": {
			PrevID:   1,
			Size:     5,
			TODOsLen: 0,
		},
		"One read with prev id = 3": {
			PrevID:   3,
			Size:     1,
			TODOsLen: 1,
		},
		"One read with prev id = 1": {
			PrevID:   1,
			Size:     1,
			TODOsLen: 0,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			var q string
			if tc.PrevID != 0 && tc.Size != 0 {
				q = fmt.Sprintf("?prev_id=%d&size=%d", tc.PrevID, tc.Size)
			} else if tc.PrevID != 0 {
				q = fmt.Sprintf("?prev_id=%d", tc.PrevID)
			} else if tc.Size != 0 {
				q = fmt.Sprintf("?size=%d", tc.Size)
			}

			resp, err := http.Get(srv.URL + "/todos" + q)
			if err != nil {
				t.Errorf("リクエストに失敗しました: %v", err)
				return
			}
			defer resp.Body.Close()

			body := map[string]interface{}{
				"todos": []map[string]interface{}{},
			}
			err = json.NewDecoder(resp.Body).Decode(&body)
			if err != nil {
				t.Error("jsonのデコードに失敗しました", err)
				return
			}

			got, ok := body["todos"].([]interface{})
			if !ok {
				t.Error("todos field が見つかりません")
				return
			}

			if len(got) != tc.TODOsLen {
				t.Errorf("期待していない todos の長さです, got = %d, want = %d", len(got), tc.TODOsLen)
			}
		})
	}
}
