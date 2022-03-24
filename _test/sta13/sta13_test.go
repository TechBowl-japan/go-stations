package sta13_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStation13(t *testing.T) {
	dbPath := "./temp_test.db"
	if err := os.Setenv("DB_PATH", dbPath); err != nil {
		t.Errorf("dbPathのセットに失敗しました。%v", err)
		return
	}

	t.Cleanup(func() {
		if err := os.Remove(dbPath); err != nil {
			t.Errorf("テスト用のDBファイルの削除に失敗しました: %v", err)
			return
		}
	})

	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		t.Errorf("データベースの作成に失敗しました: %v", err)
		return
	}

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

	_, err = stmt.Exec("todo subject")
	if err != nil {
		t.Errorf("todoの追加に失敗しました: %v", err)
		return
	}

	r := router.NewRouter(todoDB)
	srv := httptest.NewServer(r)
	defer srv.Close()

	testcases := map[string]struct {
		ID                 float64
		Subject            string
		Description        string
		WantHTTPStatusCode int
	}{
		"ID is empty": {
			WantHTTPStatusCode: http.StatusBadRequest,
		},
		"Subject is empty": {
			ID:                 1,
			WantHTTPStatusCode: http.StatusBadRequest,
		},
		"Description is empty": {
			ID:                 1,
			Subject:            "todo subject 1",
			WantHTTPStatusCode: http.StatusOK,
		},
		"Subject and Description is not empty": {
			ID:                 1,
			Subject:            "todo subject 2",
			Description:        "todo description 2",
			WantHTTPStatusCode: http.StatusOK,
		},
	}

	for name, tc := range testcases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, srv.URL+"/todos",
				bytes.NewBufferString(fmt.Sprintf(`{"id":%d,"subject":"%s","description":"%s"}`, int64(tc.ID), tc.Subject, tc.Description)))
			if err != nil {
				t.Errorf("リクエストの作成に失敗しました: %v", err)
				return
			}
			req.Header.Add("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("リクエストの送信に失敗しました: %v", err)
				return
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Errorf("レスポンスのクローズに失敗しました: %v", err)
					return
				}
			}()

			if resp.StatusCode != tc.WantHTTPStatusCode {
				t.Errorf("期待していない HTTP status code です, got = %d, want = %d", resp.StatusCode, tc.WantHTTPStatusCode)
				return
			}

			if tc.WantHTTPStatusCode != http.StatusOK {
				return
			}

			var m map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
				t.Errorf("レスポンスのデコードに失敗しました: %v", err)
				return
			}

			v, ok := m["todo"]
			if !ok {
				t.Errorf("レスポンスにtodoが含まれていません")
				return
			}

			got, ok := v.(map[string]interface{})
			if !ok {
				t.Error("todo field が Map に変換できません")
				return
			}
			want := map[string]interface{}{
				"id":          tc.ID,
				"subject":     tc.Subject,
				"description": tc.Description,
			}

			now := time.Now().UTC()
			diff := cmp.Diff(got, want, cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
				switch k {
				case "created_at", "updated_at":
					vv, ok := v.(string)
					if !ok {
						t.Errorf("日付が文字列に変換できません, got = %+v", k)
						return true
					}
					if tt, err := time.Parse(time.RFC3339, vv); err != nil {
						t.Errorf("日付が期待しているフォーマットではありません, got = %s", k)
					} else if now.Before(tt) {
						t.Errorf("日付が未来の日付になっています, got = %s", tt)
					}
					return true
				}

				return false
			}))

			if diff != "" {
				t.Error("期待していない値です\n", diff)
				return
			}
		})
	}
}
