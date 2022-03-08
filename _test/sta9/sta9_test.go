package sta9_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/osamingo/go-todo-app/db"
	"github.com/osamingo/go-todo-app/handler/router"
)

func TestStation9(t *testing.T) {
	dbPath := "./temp_test.db"
	if err := os.Setenv("DB_PATH", dbPath); err != nil {
		t.Error("dbPathのセットに失敗しました。", err)
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
		t.Error("DBの作成に失敗しました。", err)
		return
	}
	r := router.NewRouter(todoDB)
	srv := httptest.NewServer(r)
	defer srv.Close()

	testcases := map[string]struct {
		Subject            string
		Description        string
		WantHTTPStatusCode int
	}{
		"Subject is empty": {
			WantHTTPStatusCode: http.StatusBadRequest,
		},
		"Description is empty": {
			Subject:            "todo subject",
			WantHTTPStatusCode: http.StatusOK,
		},
		"Subject and Description is not empty": {
			Subject:            "todo subject",
			Description:        "todo description",
			WantHTTPStatusCode: http.StatusOK,
		},
	}

	for name, tc := range testcases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			resp, err := http.Post(srv.URL+"/todos", "application/json",
				bytes.NewBufferString(fmt.Sprintf(`{"subject":"%s","description":"%s"}`, tc.Subject, tc.Description)))
			if err != nil {
				t.Error("リクエストの送信に失敗しました。", err)
				return
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Error("レスポンスのクローズに失敗しました。", err)
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
				t.Error("レスポンスのデコードに失敗しました。", err)
				return
			}

			v, ok := m["todo"]
			if !ok {
				t.Error("レスポンスの中にtodoがありません。")
				return
			}

			got, ok := v.(map[string]interface{})
			if !ok {
				t.Error("レスポンスの中のtodoがmapではありません。")
				return
			}
			want := map[string]interface{}{
				"subject":     tc.Subject,
				"description": tc.Description,
			}

			now := time.Now().UTC()
			diff := cmp.Diff(got, want, cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
				switch k {
				case "id":
					if vv, _ := v.(float64); vv == 0 {
						t.Errorf("id を数値に変換できません, got = %s", k)
					}
					return true
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
			}
		})
	}
}
